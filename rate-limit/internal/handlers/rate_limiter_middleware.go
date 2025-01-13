package handlers

import (
	"log"
	"net"
	"net/http"
	"rate-limit/internal/core/domain"
	"strings"
)

type RateLimiterMiddleware struct {
    limiter domain.RateLimiter
}

func NewRateLimiterMiddleware(limiter domain.RateLimiter) *RateLimiterMiddleware {
    return &RateLimiterMiddleware{limiter: limiter}
}

func extractClientIP(r *http.Request) string {
    // Ordem de precedência para capturar o IP
    ipSources := []string{
        "X-Forwarded-For",
        "X-Real-IP",
        "Remote-Addr",
    }

    for _, source := range ipSources {
        ip := r.Header.Get(source)

        // Se X-Forwarded-For tiver múltiplos IPs, pega o primeiro
        if source == "X-Forwarded-For" {
            ips := strings.Split(ip, ",")
            ip = strings.TrimSpace(ips[0])
        }

        // Remover a porta se estiver presente
        if ip != "" {
            // Tentar fazer parse do IP
            host, _, err := net.SplitHostPort(ip)
            if err == nil {
                // Se conseguiu fazer split, usa o host (IP sem porta)
                ip = host
            }

            // Validar se é um IP válido
            parsedIP := net.ParseIP(ip)
            if parsedIP != nil {
                log.Printf("IP extraído de %s: %s", source, ip)
                return ip
            }
        }
    }

    // Tratamento especial para RemoteAddr
    if r.RemoteAddr != "" {
        host, _, err := net.SplitHostPort(r.RemoteAddr)
        if err == nil {
            parsedIP := net.ParseIP(host)
            if parsedIP != nil {
                log.Printf("IP extraído de RemoteAddr: %s", host)
                return host
            }
        }
    }

    // Fallback para IP desconhecido
    log.Println("Não foi possível extrair um IP válido")
    return "unknown"
}

func (m *RateLimiterMiddleware) Middleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var rateLimiterReq domain.RateLimiterRequest

        // Determinar o tipo de requisição
        token := r.Header.Get("API_KEY")
        if token != "" {
            rateLimiterReq = domain.RateLimiterRequest{
                Key:  token,
                Type: "token",
            }
        } else {
            // Extrair IP do cabeçalho
            clientIP := extractClientIP(r)

            log.Printf("Cliente IP: %s", clientIP)

            rateLimiterReq = domain.RateLimiterRequest{
                Key:  clientIP,
                Type: "ip",
            }
        }

        log.Printf("Verificando rate limit para: %+v", rateLimiterReq)

        allowed, err := m.limiter.IsAllowed(rateLimiterReq)
        if err != nil {
            log.Printf("Erro no rate limit: %v", err)
            http.Error(w,
                "Você atingiu o número máximo de solicitações permitidas",
                http.StatusTooManyRequests,
            )
            return
        }

        if !allowed {
            log.Println("Requisição não permitida")
            http.Error(w,
                "Você atingiu o número máximo de solicitações permitidas",
                http.StatusTooManyRequests,
            )
            return
        }

        // Se chegou aqui, a requisição é permitida
        log.Println("Requisição permitida, prosseguindo...")
        next.ServeHTTP(w, r)
    }
}