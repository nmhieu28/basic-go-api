package cache

import (
	configs "backend/pkg/config"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisClient(appConfig *configs.AppConfig) Cache {
	config := appConfig.Redis
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	options := &redis.Options{
		Addr:     addr,
		DB:       0,
		Username: config.UserName,
		Password: config.Password,
		PoolSize: 10,
	}

	client := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {

	}

	return &RedisCache{
		client: client,
	}
}

func (r *RedisCache) Connect(ctx context.Context) error {
	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return nil
}

func (r *RedisCache) SelectDatabase(ctx context.Context, db int) error {

	if err := r.client.Do(ctx, "SELECT", db).Err(); err != nil {
		return fmt.Errorf("failed to select database %d: %w", db, err)
	}
	return nil
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("error getting key %s: %w", key, err)
	}
	return val, nil
}

// ttl: seconds
func (r *RedisCache) Set(ctx context.Context, key string, value string, ttl int) error {
	if err := r.client.Set(ctx, key, value, time.Duration(ttl)*time.Second).Err(); err != nil {
		return fmt.Errorf("error setting key %s: %w", key, err)
	}
	return nil
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("error deleting key %s: %w", key, err)
	}
	return nil
}

func (r *RedisCache) Disconnect() error {
	if err := r.client.Close(); err != nil {
		return fmt.Errorf("error disconnecting from Redis: %w", err)
	}
	fmt.Println("Disconnected from Redis")
	return nil
}
