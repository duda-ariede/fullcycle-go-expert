package usecases

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/viper"

	"rate-limit/internal/core/domain"
	"rate-limit/internal/ports"
)

type RateLimiterService struct {
    repository ports.RateLimiterRepository
    config     domain.RateLimiterConfig
    mu         sync.Mutex
}

func NewRateLimiterService(
    repo ports.RateLimiterRepository,
    config domain.RateLimiterConfig,
) domain.RateLimiter {
    return &RateLimiterService{
        repository: repo,
        config:     config,
    }
}

func (s *RateLimiterService) loadLimitsFromConfig() {
    s.config = domain.RateLimiterConfig{
        Enabled: viper.GetBool("RATE_LIMIT_ENABLED"),
        Limits: map[string]domain.LimitConfig{
            "ip": {
                MaxRequests:      viper.GetInt("RATE_LIMIT_IP_MAX_REQUESTS"),
                BlockDurationMin: viper.GetInt("RATE_LIMIT_IP_BLOCK_DURATION"),
            },
            "token": {
                MaxRequests:      viper.GetInt("RATE_LIMIT_TOKEN_MAX_REQUESTS"),
                BlockDurationMin: viper.GetInt("RATE_LIMIT_TOKEN_BLOCK_DURATION"),
            },
        },
    }
}

// internal/core/usecases/rate_limiter_service.go
func (s *RateLimiterService) IsAllowed(req domain.RateLimiterRequest) (bool, error) {
    // Recarregar configurações
    s.loadLimitsFromConfig()

    // Log de depuração
    log.Printf("Rate Limit Check - Key: %s, Type: %s", req.Key, req.Type)

    // Verificar se o rate limit está habilitado
    if !s.config.Enabled {
        log.Println("Rate limit is disabled")
        return true, nil
    }

    // Determinar o limite de acordo com o tipo de requisição
    limitConfig, exists := s.config.Limits[req.Type]
    if !exists {
        log.Printf("No limit configured for type: %s", req.Type)
        return true, nil
    }

    // Chaves para rastreamento
    windowKey := fmt.Sprintf("rate_limit:%s:%s:window", req.Type, req.Key)
    blockedKey := fmt.Sprintf("rate_limit:%s:%s:blocked", req.Type, req.Key)

    // Verificar se está bloqueado
    blockedUntil, err := s.repository.Get(blockedKey)
    if err == nil && blockedUntil != "" {
        blockedTime, _ := strconv.ParseInt(blockedUntil, 10, 64)
        currentTime := time.Now().Unix()

        log.Printf("Blocked until: %d, Current time: %d", blockedTime, currentTime)

        if currentTime < blockedTime {
            log.Printf("Key is still blocked: %s. Time remaining: %d seconds",
                req.Key, blockedTime - currentTime)
            return false, fmt.Errorf("key is blocked")
        } else {
            // Tempo de bloqueio expirou, remover chave de bloqueio
            log.Printf("Bloqueio expirado para a chave: %s", req.Key)
            err := s.repository.Delete(blockedKey)
            if err != nil {
                log.Printf("Erro ao remover chave de bloqueio: %v", err)
            }

            // Também limpar a janela de contagem
            err = s.repository.Delete(windowKey)
            if err != nil {
                log.Printf("Erro ao remover janela de contagem: %v", err)
            }
        }
    }

    // Incrementar contador de requisições
    currentCount, err := s.repository.Increment(windowKey)
    if err != nil {
        log.Printf("Error incrementing counter: %v", err)
        return false, err
    }

    // Definir expiração para a janela de contagem (1 minuto)
    err = s.repository.Set(windowKey, currentCount, 1)
    if err != nil {
        log.Printf("Error setting window key: %v", err)
        return false, err
    }

    // Verificar se o limite de requisições foi excedido
    if currentCount > int64(limitConfig.MaxRequests) {
        // Bloquear a chave pelo tempo configurado
        blockUntil := time.Now().Add(time.Duration(limitConfig.BlockDurationMin) * time.Minute).Unix()
        err := s.repository.Set(blockedKey, blockUntil, limitConfig.BlockDurationMin)
        if err != nil {
            log.Printf("Error blocking key: %v", err)
            return false, err
        }
        log.Printf("Rate limit exceeded for key: %s. Blocked until %d", req.Key, blockUntil)
        return false, fmt.Errorf("rate limit exceeded")
    }

    log.Println("Request allowed")
    return true, nil
}

func (s *RateLimiterService) BlockKey(key string, duration int) error {
    blockedKey := fmt.Sprintf("rate_limit:manual:%s:blocked", key)
    blockUntil := time.Now().Add(time.Duration(duration) * time.Minute).Unix()
    return s.repository.Set(blockedKey, blockUntil, duration)
}

// Método para atualizar a configuração de forma thread-safe
func (s *RateLimiterService) UpdateConfig(newConfig domain.RateLimiterConfig) {
    s.mu.Lock()
    defer s.mu.Unlock()

    log.Printf("Atualizando configuração de Rate Limiter: Enabled=%v", newConfig.Enabled)

    // Imprimir detalhes dos limites para log
    for key, limit := range newConfig.Limits {
        log.Printf("Limite para %s: MaxRequests=%d, BlockDuration=%d min",
            key, limit.MaxRequests, limit.BlockDurationMin)
    }

    s.config = newConfig
}