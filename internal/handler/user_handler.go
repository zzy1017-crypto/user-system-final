package handler

import (
	"strconv"
	"strings"

	"user-system-final/internal/service"

	"github.com/gin-gonic/gin"
)

// UserHandler 是一个处理用户相关HTTP请求的处理器结构体，包含一个UserService字段，用于处理用户相关的业务逻辑
type UserHandler struct {
	service *service.UserService
}

// NewUserHandler 创建一个新的 UserHandler 实例，接受一个 UserService 作为参数，并返回一个 UserHandler 对象，提供给路由处理函数使用以处理用户相关的请求
func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// Login 处理用户登录请求，接受一个包含用户名和密码的JSON请求体，调用 UserService 的 Login 方法进行用户认证，并返回一个JWT令牌或错误信息
func (h *UserHandler) Login(c *gin.Context) {

	//定义一个匿名结构体来绑定请求体中的JSON数据，包含用户名和密码字段
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	c.ShouldBindJSON(&req) //将请求体中的JSON数据绑定到req结构体中，如果绑定失败则返回400错误响应
	if req.Username == "" || req.Password == "" {
		c.JSON(400, gin.H{"error": "username and password are required"})
		return
	}

	//调用 UserService 的 Login 方法进行用户认证，传入用户名和密码，如果认证成功则返回一个JWT令牌，否则返回错误信息
	token, err := h.service.Login(req.Username, req.Password)

	//如果登录失败，返回400错误响应，并包含错误信息
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": token}) //如果登录成功，返回200响应，并包含生成的JWT令牌
}

// Logout 处理用户注销请求，从请求头中获取JWT令牌，调用 UserService 的 Logout 方法进行用户注销，并返回注销结果
func (h *UserHandler) Logout(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")         //从请求头中获取Authorization字段，通常格式为"Bearer <token>"
	token := strings.TrimPrefix(authHeader, "Bearer ") //去掉"Bearer "前缀，获取纯粹的token字符串

	//调用 UserService 的 Logout 方法进行用户注销，传入JWT令牌，如果注销成功则返回成功消息，否则返回错误信息
	err := h.service.Logout(token)
	//如果注销失败，返回500错误响应，并包含错误信息
	if err != nil {
		c.JSON(500, gin.H{"error": "logout failed"})
		return
	}

	c.JSON(200, gin.H{"msg": "logout success"}) //如果注销成功，返回200响应，并包含成功消息
}

// Profile 处理获取用户个人资料的请求，从上下文中获取用户ID，调用 UserService 的 GetUserProfile 方法获取用户信息，并返回用户资料或错误信息
func (h *UserHandler) Profile(c *gin.Context) {

	userIDVal, _ := c.Get("userID") //从上下文中获取用户ID，之前在AuthMiddleware中已经将用户ID存储在上下文中，这里直接获取并转换为字符串类型

	userIDStr := userIDVal.(string) //将获取到的用户ID值断言为字符串类型，如果断言失败会引发panic，这里假设之前的中间件已经确保了这个值是字符串类型

	userID, _ := strconv.Atoi(userIDStr) //将用户ID字符串转换为整数类型，如果转换失败会引发panic，这里假设用户ID是一个有效的整数字符串

	//调用 UserService 的 GetUserProfile 方法获取用户信息，传入用户ID，如果获取成功则返回用户资料，否则返回错误信息
	user, err := h.service.GetUserProfile(userID)
	//如果获取用户资料失败，返回500错误响应，并包含错误信息
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to get user profile",
		})
		return
	}

	c.JSON(200, user) //如果获取用户资料成功，返回200响应，并包含用户信息作为JSON数据
}
