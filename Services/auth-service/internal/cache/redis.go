package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// concrete implementation
type redisClient struct {
	client *redis.Client
}

// Constructor – ONLY responsibility: establish connection
func NewRedisClient(host, port string) (RedisClient, error) {
	addr := fmt.Sprintf("%s:%s", host, port)

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
	})

	// verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &redisClient{client: rdb}, nil
}

/*
INTERFACE — unchanged (as requested)
*/
type RedisClient interface {
	Set(ctx context.Context, key string, val string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	SetToken(ctx context.Context, key string, val string, ttl time.Duration) error
	Close() error
}

/*
	METHOD IMPLEMENTATIONS — unchanged signatures
*/

func (r *redisClient) Set(ctx context.Context, key string, val string, ttl time.Duration) error {
	return r.client.Set(ctx, key, val, ttl).Err()
}

func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *redisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Semantic wrapper — still uses same underlying Set
func (r *redisClient) SetToken(ctx context.Context, key string, val string, ttl time.Duration) error {
	return r.Set(ctx, key, val, ttl)
}
func (r *redisClient) Close() error {
	if r.client != nil {
		r.client.Close()
	}
	return nil
}
