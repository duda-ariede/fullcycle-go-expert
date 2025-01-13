package domain

type RateLimiterConfig struct {
    Enabled bool
    Limits  map[string]LimitConfig
}

type LimitConfig struct {
    MaxRequests      int
    BlockDurationMin int
}

type RateLimiterRequest struct {
    Key   string
    Type  string // "ip" ou "token"
}

type RateLimiter interface {
    IsAllowed(req RateLimiterRequest) (bool, error)
    BlockKey(key string, duration int) error
}