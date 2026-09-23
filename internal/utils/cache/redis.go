package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"movie-ticket/config"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

var RedisService *RedisClient

func InitRedis(cfg *config.Config) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: "", // no password set by default
		DB:       0,  // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Printf("Warning: Unable to connect to Redis at %s: %v. Caching will be bypassed gracefully.", cfg.RedisURL, err)
	} else {
		log.Println("Successfully connected to Redis cache.")
	}

	RedisService = &RedisClient{client: rdb}
	return RedisService
}

func (r *RedisClient) Client() *redis.Client {
	return r.client
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if r == nil || r.client == nil {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
	if r == nil || r.client == nil {
		return redis.Nil
	}
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisClient) Delete(ctx context.Context, keys ...string) error {
	if r == nil || r.client == nil || len(keys) == 0 {
		return nil
	}
	return r.client.Del(ctx, keys...).Err()
}

func (r *RedisClient) DeleteByPrefix(ctx context.Context, prefix string) error {
	if r == nil || r.client == nil {
		return nil
	}
	var cursor uint64
	for {
		keys, nextCursor, err := r.client.Scan(ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			_ = r.client.Del(ctx, keys...).Err()
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}
