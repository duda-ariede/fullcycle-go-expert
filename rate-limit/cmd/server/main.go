package main

import (
	"log"
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"rate-limit/internal/adapters/redis"
	"rate-limit/internal/core/domain"
	"rate-limit/internal/core/usecases"
	"rate-limit/internal/handlers"
	"rate-limit/pkg/config"
)

func main() {
    // Carregar configurações
    config.LoadConfig()

    // Criar repositório Redis
    redisStore := redis.NewRedisStore(
        viper.GetString("REDIS_HOST"),
        viper.GetString("REDIS_PORT"),
    )

    // Função para atualizar a configuração do rate limiter
    updateRateLimiterConfig := func() domain.RateLimiterConfig {
        ipMaxRequests, ipBlockDuration,
        tokenMaxRequests, tokenBlockDuration := config.GetRateLimiterConfig()

        return domain.RateLimiterConfig{
            Enabled: viper.GetBool("RATE_LIMIT_ENABLED"),
            Limits: map[string]domain.LimitConfig{
                "ip": {
                    MaxRequests:      ipMaxRequests,
                    BlockDurationMin: ipBlockDuration,
                },
                "token": {
                    MaxRequests:      tokenMaxRequests,
                    BlockDurationMin: tokenBlockDuration,
                },
            },
        }
    }

    // Criar serviço de rate limiter inicial
    rateLimiterService := usecases.NewRateLimiterService(
        redisStore,
        updateRateLimiterConfig(),
    )

    // Criar middleware de rate limiter
    middleware := handlers.NewRateLimiterMiddleware(rateLimiterService)

    // Goroutine para monitorar alterações de configuração
    go func() {
        // Canal do Viper para mudanças de configuração
        configChanges := make(chan struct{}, 1)

        viper.OnConfigChange(func(e fsnotify.Event) {
            select {
            case configChanges <- struct{}{}:
            default:
            }
        })

        for range configChanges {
            log.Println("Configuração alterada, atualizando rate limiter...")

            // Atualizar configuração do rate limiter
            newConfig := updateRateLimiterConfig()

            // Type assertion para acessar o método UpdateConfig
            if service, ok := rateLimiterService.(*usecases.RateLimiterService); ok {
                service.UpdateConfig(newConfig)
            } else {
                log.Println("Não foi possível atualizar a configuração do rate limiter")
            }
        }
    }()

    // Exemplo de rota protegida com rate limiting
    http.HandleFunc("/", middleware.Middleware(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Request successful"))
    }))

    // Iniciar servidor
    log.Println("Servidor iniciado na porta 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Erro ao iniciar o servidor: %v", err)
    }
}