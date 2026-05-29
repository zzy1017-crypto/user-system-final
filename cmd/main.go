package main

import (
	"github.com/gin-gonic/gin"

	"user-system-final/internal/auth"
	"user-system-final/internal/cache"
	"user-system-final/internal/handler"
	"user-system-final/internal/repository"
	"user-system-final/internal/service"
)

func main() {

	r := gin.Default()

	// 初始化
	repo := repository.NewMemoryRepo()
	redisClient := cache.NewRedisClient("localhost:6379")
	svc := service.NewUserService(repo, redisClient)
	h := handler.NewUserHandler(svc)

	// 路由
	r.POST("/login", h.Login)

	authGroup := r.Group("/")
	authGroup.Use(
		auth.AuthMiddleware(redisClient),
		auth.RateLimitMiddleware(redisClient, 5),
	)
	authGroup.POST("/logout", h.Logout)
	authGroup.GET("/profile", h.Profile)

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Run(":8080")
}
