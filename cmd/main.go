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

	r := gin.Default() //开启Gin框架的默认设置，创建一个新的Gin引擎实例，准备处理HTTP请求

	repo := repository.NewMemoryRepo()                    //创建一个内存中的用户仓库实例，提供用户数据的存储和访问功能
	redisClient := cache.NewRedisClient("localhost:6379") //创建一个Redis客户端实例，连接到本地的Redis服务器，提供缓存和数据存储功能
	svc := service.NewUserService(repo, redisClient)      //创建一个用户服务实例，接受用户仓库和Redis客户端作为依赖，提供用户相关的业务逻辑处理功能
	h := handler.NewUserHandler(svc)                      //创建一个用户处理器实例，接受用户服务作为依赖，提供处理用户相关HTTP请求的功能

	r.POST("/login", h.Login) //定义一个POST路由，路径为/login，处理函数为h.Login，用于处理用户登录请求

	authGroup := r.Group("/") //创建一个新的路由组，路径前缀为"/"，用于处理需要认证和授权的用户相关请求

	//在authGroup路由组中使用认证和限流中间件，传入Redis客户端作为参数，确保只有经过认证的用户才能访问这些路由，并且限制每个用户的请求频率
	authGroup.Use(
		auth.AuthMiddleware(redisClient),
		auth.RateLimitMiddleware(redisClient, 5),
	)
	authGroup.POST("/logout", h.Logout)  //在authGroup路由组中定义一个POST路由，路径为/logout，处理函数为h.Logout，用于处理用户注销请求
	authGroup.GET("/profile", h.Profile) //在authGroup路由组中定义一个GET路由，路径为/profile，处理函数为h.Profile，用于处理获取用户个人资料的请求

	r.Use(gin.Logger())   //使用Gin框架的日志中间件，记录HTTP请求的相关信息，例如请求方法、路径、状态码等，方便调试和监控
	r.Use(gin.Recovery()) //使用Gin框架的恢复中间件，捕获运行时的panic，防止程序崩溃，并返回500错误响应

	r.Run(":8080") //启动Gin服务器，监听8080端口，准备接受和处理HTTP请求
}
