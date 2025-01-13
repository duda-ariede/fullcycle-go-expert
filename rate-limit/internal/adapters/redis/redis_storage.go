package redis

import (
	"context"
	"time"

	"rate-limit/internal/ports"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
    client *redis.Client
}

func NewRedisStore(host, port string) ports.RateLimiterRepository {
    rdb := redis.NewClient(&redis.Options{
        Addr: host + ":" + port,
    })

    return &RedisStore{client: rdb}
}

func (r *RedisStore) Increment(key string) (int64, error) {
    ctx := context.Background()
    return r.client.Incr(ctx, key).Result()
}

func (r *RedisStore) Set(key string, value interface{}, expiration int) error {
    ctx := context.Background()
    return r.client.Set(ctx, key, value, time.Duration(expiration)*time.Minute).Err()
}

func (r *RedisStore) Get(key string) (string, error) {
    ctx := context.Background()
    return r.client.Get(ctx, key).Result()
}

func (r *RedisStore) Delete(key string) error {
    ctx := context.Background()
    return r.client.Del(ctx, key).Err()
}