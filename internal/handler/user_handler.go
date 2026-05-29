package handler

import (
	"strconv"
	"strings"

	"user-system-final/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) Login(c *gin.Context) {

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	c.ShouldBindJSON(&req)

	token, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func (h *UserHandler) Logout(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	err := h.service.Logout(token)
	if err != nil {
		c.JSON(500, gin.H{"error": "logout failed"})
		return
	}

	c.JSON(200, gin.H{"msg": "logout success"})
}

func (h *UserHandler) Profile(c *gin.Context) {

	userIDVal, _ := c.Get("userID")

	userIDStr := userIDVal.(string)

	userID, _ := strconv.Atoi(userIDStr)

	user, err := h.service.GetUserProfile(userID)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to get user profile",
		})
		return
	}

	c.JSON(200, user)
}
