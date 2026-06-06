package repository

import "user-system-final/internal/model"

// UserRepository 定义了一个用户仓库接口，包含一个方法GetUserByName，用于根据用户名获取用户信息，返回一个用户对象和错误信息
type UserRepository interface {
	GetUserByName(name string) (*model.User, error)
}
