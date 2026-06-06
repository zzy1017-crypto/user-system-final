package auth

import (
	"net/http"
	"strings"

	"user-system-final/internal/cache"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 是一个JWT认证中间件，验证请求中的JWT令牌是否合法，并将用户ID存储在上下文中供后续处理使用
func AuthMiddleware(rdb *cache.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization") //从请求头中获取Authorization字段，通常格式为"Bearer <token

		//如果Authorization字段为空，说明没有提供token，返回 401 Unauthorized 错误响应，并中止请求处理
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no token"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ") //去掉"Bearer "前缀，获取纯粹的token字符串

		//解析token，验证其有效性，并提取claims（载荷）
		_, err := ParseToken(tokenStr)

		//如果解析token失败，说明token无效，返回 401 Unauthorized 错误响应，并中止请求处理
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		key := "login:token:" + tokenStr //构造一个Redis键，格式为"login:token:<token字符串>"，用于在Redis中存储和查询用户登录状态

		//从Redis中获取与该token对应的用户ID，如果获取失败，说明token已过期或无效，返回 401 Unauthorized 错误响应，并中止请求处理
		userID, err := rdb.Get(key)
		//如果从Redis获取用户ID失败，说明token已过期或无效，返回 401 Unauthorized 错误响应，并中止请求处理
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			c.Abort()
			return
		}

		c.Set("userID", userID) //将用户ID存储在上下文中，供后续处理使用，例如在后续的handler中可以通过c.Get("userID")来获取当前请求的用户ID

		c.Next() //放行请求，继续处理后续的handler
	}

}
