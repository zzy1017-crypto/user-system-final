package cache

import (
	"context"
	"encoding/json"
	"time"
	"user-system-final/internal/logger"

	"github.com/redis/go-redis/v9"
)

// RedisClient 包装了 go-redis 的客户端对象，提供了一个结构体来封装 Redis 客户端，方便在程序中使用和扩展
type RedisClient struct {
	client *redis.Client   //Redis 客户端对象，用于执行 Redis 命令和操作，封装了 go-redis 库的客户端功能，可以在程序中其他地方使用这个对象来操作 Redis
	ctx    context.Context //全局上下文对象，go-redis库的操作需要一个上下文参数，这里创建一个全局的背景上下文，可以在程序中其他地方使用cache.Ctx来传递这个上下文
}

// NewRedisClient 创建一个新的 Redis 客户端，接受 Redis 的连接信息作为参数，返回一个 RedisClient 对象
func NewRedisClient(addr string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr, //Redis 服务器的地址，格式为 "host:port"，例如 "localhost:6379"，用于连接到 Redis 服务器
	})

	return &RedisClient{
		client: rdb,                  //返回一个 RedisClient 对象，包含了创建的 Redis 客户端，可以在程序中其他地方使用这个对象来操作 Redis
		ctx:    context.Background(), //创建一个背景上下文，表示没有特定的取消或超时机制，可以在程序中其他地方使用这个上下文来传递给 Redis 操作
	}
}

// Set 设置一个键值对到 Redis 中，接受一个键、一个值和一个过期时间作为参数，返回错误信息
func (r *RedisClient) Set(key string, value string, exp time.Duration) error {
	err := r.client.Set(r.ctx, key, value, exp).Err()
	if err != nil {
		logger.Error("Redis SET failed:" + key)
	} else {
		logger.Info("Redis SET:" + key)
	}
	return err
}

// Get 从 Redis 中获取一个键的值，接受一个键作为参数，返回值和错误信息
func (r *RedisClient) Get(key string) (string, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		logger.Error("Redis GET failed:" + key)
	} else {
		logger.Info("Redis GET :" + key)
	}
	return val, err
}

// Del 从 Redis 中删除一个键，接受一个键作为参数，返回错误信息
func (r *RedisClient) Del(key string) error {
	err := r.client.Del(r.ctx, key).Err()
	if err != nil {
		logger.Error("Redis DEL failed:" + key)
	} else {
		logger.Info("Redis DEL:" + key)
	}
	return err
}

// SetJSON 将一个结构体或对象以 JSON 格式存储到 Redis 中，接受一个键、一个值和一个过期时间作为参数，返回错误信息
func (r *RedisClient) SetJSON(key string, value interface{}, exp time.Duration) error {

	//将值序列化为 JSON 格式的字节切片，如果序列化失败，返回错误信息
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(r.ctx, key, data, exp).Err() //将序列化后的 JSON 数据存储到 Redis 中，使用 Set 命令，传递键、值和过期时间，如果操作失败，返回错误信息
}

// GetJSON 从 Redis 中获取一个键的值，并将其反序列化为指定的结构体或对象，接受一个键和目标对象作为参数，返回错误信息
func (r *RedisClient) GetJSON(key string, dest interface{}) error {

	//从 Redis 中获取一个键的值，使用 Get 命令，如果操作失败，返回错误信息
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest) //将获取到的值反序列化为目标对象，使用 json.Unmarshal 函数
}

// Incr 对指定的键进行自增操作，返回自增后的值和错误信息
func (r *RedisClient) Incr(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

// Expire 设置指定键的过期时间，接受一个键和过期时间作为参数，返回错误信息
func (r *RedisClient) Expire(key string, exp time.Duration) error {
	return r.client.Expire(r.ctx, key, exp).Err()
}
