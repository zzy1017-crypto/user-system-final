package auth

import (
	"fmt"
	"net/http"
	"time"

	"user-system-final/internal/cache"
	"user-system-final/internal/logger"

	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(rdb *cache.RedisClient, limit int) gin.HandlerFunc {

	return func(c *gin.Context) {

		userIDVal, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "no user",
			})

			c.Abort()
			return
		}

		userID := userIDVal.(string)

		now := time.Now().Format("200601021504")
		key := fmt.Sprintf("rate_limit:%s:%s", userID, now)

		count, err := rdb.Incr(key)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "redis error",
			})
			c.Abort()
			return
		}

		if count == 1 {
			_ = rdb.Expire(key, time.Minute)
		}

		if count > int64(limit) {
			logger.Error("rate limit exceeded for user:" + userID)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
