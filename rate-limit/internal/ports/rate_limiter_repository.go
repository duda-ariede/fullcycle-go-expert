package ports

type RateLimiterRepository interface {
    Increment(key string) (int64, error)
    Set(key string, value interface{}, expiration int) error
    Get(key string) (string, error)
    Delete(key string) error
}