package auth

import (
	"fmt"
	"net/http"
	"time"

	"user-system-final/internal/cache"
	"user-system-final/internal/logger"

	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware 是一个限流中间件，限制每个用户在一定时间内的请求次数，防止恶意请求和过载
func RateLimitMiddleware(rdb *cache.RedisClient, limit int) gin.HandlerFunc {

	return func(c *gin.Context) {

		//从上下文中获取用户ID，如果不存在，说明用户未认证，返回 401 Unauthorized 错误响应，并中止请求处理
		userIDVal, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "no user",
			})

			c.Abort()
			return
		}

		userID := userIDVal.(string) //将用户ID转换为字符串类型，准备进行限流计数

		now := time.Now().Format("200601021504")            //获取当前时间，格式化为"yyyyMMddHHmm"的字符串形式，作为限流的时间窗口
		key := fmt.Sprintf("rate_limit:%s:%s", userID, now) //构造一个Redis键，格式为"rate_limit:<userID>:<当前时间窗口>"，用于在Redis中存储和查询用户的请求计数

		//使用 Redis 的 INCR 命令增加请求计数器，count是当前计数值，err是操作是否成功的错误信息
		count, err := rdb.Incr(key)
		//如果 Redis 操作失败，返回 500 Internal Server Error 错误响应，并中止请求处理
		if err != nil {
			c.JSON(500, gin.H{
				"error": "redis error",
			})
			c.Abort()
			return
		}

		//如果计数器的值为 1，说明这是第一次请求，需要设置过期时间为 1 分钟，确保计数器在 1 分钟后自动重置
		if count == 1 {
			_ = rdb.Expire(key, time.Minute) //设置Redis键的过期时间为1分钟，确保计数器在1分钟后自动重置，允许用户在下一分钟继续请求
		}

		//如果请求数超过限制，记录日志并返回 429 Too Many Requests 错误响应，并中止请求处理
		if count > int64(limit) {
			logger.Error("rate limit exceeded for user:" + userID)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next() //放行请求，继续处理后续的 handler
	}
}
