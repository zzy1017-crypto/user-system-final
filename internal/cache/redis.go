package cache

import (
	"context"
	"encoding/json"
	"time"
	"user-system-final/internal/logger"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient(addr string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisClient{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (r *RedisClient) Set(key string, value string, exp time.Duration) error {
	err := r.client.Set(r.ctx, key, value, exp).Err()
	if err != nil {
		logger.Error("Redis SET failed:" + key)
	} else {
		logger.Info("Redis SET:" + key)
	}
	return err
}

func (r *RedisClient) Get(key string) (string, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		logger.Error("Redis GET failed:" + key)
	} else {
		logger.Info("Redis GET :" + key)
	}
	return val, err
}

func (r *RedisClient) Del(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisClient) SetJSON(key string, value interface{}, exp time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(r.ctx, key, data, exp).Err()
}

func (r *RedisClient) GetJSON(key string, dest interface{}) error {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisClient) Incr(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

func (r *RedisClient) Expire(key string, exp time.Duration) error {
	return r.client.Expire(r.ctx, key, exp).Err()
}
