// Package cache provides the Redis access that the services need to keep the
// server side session state.
//
// The auth-service writes the session state and the API gateway reads it, so
// the client lives in the shared module to guarantee that both sides talk to
// the same Redis instance with the same key format.
package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrKeyNotFound is returned when a key does not exist.
//
// It is separate from the transport errors because the two mean different
// things for a caller: a missing key is a definite answer, while a transport
// error only means that the answer could not be retrieved.
var ErrKeyNotFound = errors.New("key not found")

// Client is the small set of Redis operations that the services need.
type Client interface {
	// Get returns the value of a key.
	//
	// It returns ErrKeyNotFound when the key does not exist.
	Get(ctx context.Context, key string) (string, error)
	// Set stores a value with an optional time to live.
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	// Delete removes a key.
	Delete(ctx context.Context, key string) error
	// Close releases the connections.
	Close() error
}

type redisClient struct {
	client *redis.Client
}

// NewRedisClient connects to Redis and verifies that the connection works.
func NewRedisClient(host, port string) (Client, error) {
	if host == "" || port == "" {
		return nil, errors.New("redis host and port are required")
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("connect to redis at %s: %w", addr, err)
	}

	return &redisClient{client: rdb}, nil
}

// Get implements [Client].
func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrKeyNotFound
	}
	if err != nil {
		return "", err
	}
	return value, nil
}

// Set implements [Client].
func (r *redisClient) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Delete implements [Client].
func (r *redisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Close implements [Client].
func (r *redisClient) Close() error {
	if r.client == nil {
		return nil
	}
	return r.client.Close()
}
