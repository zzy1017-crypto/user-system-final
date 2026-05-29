package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/sync/singleflight"

	"user-system-final/internal/auth"
	"user-system-final/internal/cache"
	"user-system-final/internal/logger"
	"user-system-final/internal/model"
	"user-system-final/internal/repository"
)

type UserService struct {
	repo  repository.UserRepository
	cache *cache.RedisClient

	sf singleflight.Group
}

func NewUserService(repo repository.UserRepository, cache *cache.RedisClient) *UserService {
	return &UserService{
		repo:  repo,
		cache: cache,
	}
}

func (s *UserService) Login(username, password string) (string, error) {

	logger.Info("user login attempt:" + username)

	user, err := s.repo.GetUserByName(username)
	if err != nil || user == nil {
		logger.Error("login failed:user not found:" + username)
		return "", errors.New("user not found")
	}

	if user.Password != password {
		logger.Error("login failed:wrong password:" + username)
		return "", errors.New("wrong password")
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		logger.Error("login failed:generate token error")
		return "", err
	}

	key := "login:token:" + token
	err = s.cache.Set(key, fmt.Sprintf("%d", user.ID), 24*time.Hour)
	if err != nil {
		logger.Error("login failed:redis error")
		return "", err
	}

	logger.Info("login success:" + username)

	return token, nil
}

func (s *UserService) Logout(token string) error {
	key := "login:token:" + token
	return s.cache.Del(key)
}

func (s *UserService) GetUserProfile(userID int) (*model.User, error) {

	key := fmt.Sprintf("user:profile:%d", userID)

	var user model.User

	logger.Info("start get your profile")

	err := s.cache.GetJSON(key, &user)
	if err == nil {
		logger.Info("cache hit:" + key)

		return &user, nil
	}

	logger.Info("cache miss:" + key)

	v, err, _ := s.sf.Do(key, func() (interface{}, error) {

		var u model.User

		logger.Info("query DB for user")
		err := s.cache.GetJSON(key, &u)
		if err == nil {
			logger.Info("double check hit cache")

			return &u, nil
		}
		logger.Info("double check cache miss")

		userPtr, err := s.repo.GetUserByName("test")
		if err != nil || userPtr == nil {
			return nil, err
		}

		exp := time.Hour + time.Duration(rand.Intn(300))*time.Second

		_ = s.cache.SetJSON(key, userPtr, exp)

		return userPtr, nil
	})

	if err != nil {
		return nil, err
	}

	return v.(*model.User), nil
}
