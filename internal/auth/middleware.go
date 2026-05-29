package auth

import (
	"net/http"
	"strings"

	"user-system-final/internal/cache"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(rdb *cache.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no token"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		_, err := ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		key := "login:token:" + tokenStr
		userID, err := rdb.Get(key)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
