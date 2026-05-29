package repository

import "user-system-final/internal/model"

type UserRepository interface {
	GetUserByName(name string) (*model.User, error)
}
