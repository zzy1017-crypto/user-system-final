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

// UserService 定义了一个用户服务对象，包含一个用户仓库和一个Redis客户端，以及一个singleflight.Group用于处理缓存击穿问题
type UserService struct {
	repo  repository.UserRepository
	cache *cache.RedisClient

	sf singleflight.Group
}

// NewUserService 创建一个新的 UserService 实例，接受一个用户仓库和一个Redis客户端作为参数，返回一个 UserService 对象
func NewUserService(repo repository.UserRepository, cache *cache.RedisClient) *UserService {
	return &UserService{
		repo:  repo,
		cache: cache,
	}
}

// Login 处理用户登录逻辑，接受用户名和密码作为参数，返回一个登录令牌和错误信息
func (s *UserService) Login(username, password string) (string, error) {

	logger.Info("user login attempt:" + username)

	user, err := s.repo.GetUserByName(username)

	//验证用户是否存在，如果用户不存在则返回错误信息，并记录登录失败的日志
	if err != nil || user == nil {
		logger.Error("login failed:user not found:" + username)
		return "", errors.New("user not found")
	}

	//验证用户密码是否正确，如果密码错误则返回错误信息，并记录登录失败的日志
	if user.Password != password {
		logger.Error("login failed:wrong password:" + username)
		return "", errors.New("wrong password")
	}

	//生成一个登录令牌，如果生成令牌失败则返回错误信息，并记录登录失败的日志
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		logger.Error("login failed:generate token error")
		return "", err
	}

	//将登录令牌存储到Redis中，设置一个过期时间
	key := "login:token:" + token
	err = s.cache.Set(key, fmt.Sprintf("%d", user.ID), 24*time.Hour)
	//如果存储令牌失败则返回错误信息，并记录登录失败的日志
	if err != nil {
		logger.Error("login failed:redis error")
		return "", err
	}

	logger.Info("login success:" + username) //记录登录成功的日志，包含用户名信息

	return token, nil //返回生成的登录令牌和nil表示没有错误，供调用者使用这个令牌进行后续的认证和授权操作
}

// Logout 处理用户注销逻辑，接受一个登录令牌作为参数，返回错误信息
func (s *UserService) Logout(token string) error {
	key := "login:token:" + token //构造Redis中的键，格式为 "login:token:{token}"，用于存储和管理用户的登录令牌
	err := s.cache.Del(key)
	//如果删除令牌失败则返回错误信息，并记录注销失败的日志
	if err != nil {
		logger.Error("logout failed:redis error")
		return err
	}
	//记录注销成功的日志，包含令牌信息
	logger.Info("logout success:" + token)
	return nil
}

// GetUserProfile 处理获取用户个人信息的逻辑，接受一个用户ID作为参数，返回一个用户对象和错误信息
func (s *UserService) GetUserProfile(userID int) (*model.User, error) {

	key := fmt.Sprintf("user:profile:%d", userID) //构造Redis中的键，格式为 "user:profile:{userID}"，用于存储和管理用户的个人信息缓存

	var user model.User //定义一个用户对象，用于保存从缓存或数据库中获取的用户信息

	//记录开始获取用户个人信息的日志，包含用户ID信息
	logger.Info("start get your profile")

	//尝试从Redis缓存中获取用户个人信息，如果获取成功则返回用户对象和nil表示没有错误，并记录缓存命中的日志
	err := s.cache.GetJSON(key, &user)
	if err == nil {
		logger.Info("cache hit:" + key)

		return &user, nil
	}

	logger.Info("cache miss:" + key) //记录缓存未命中的日志，包含键信息

	//使用singleflight.Group来处理缓存击穿问题，确保在高并发情况下只有一个请求会查询数据库并更新缓存，其他请求会等待这个请求完成后直接从缓存中获取数据
	v, err, _ := s.sf.Do(key, func() (interface{}, error) {

		var u model.User //定义一个用户对象，用于保存从数据库中获取的用户信息

		logger.Info("query DB for user") //记录查询数据库的日志，包含用户ID信息

		//在查询数据库之前再次尝试从Redis缓存中获取用户个人信息，以防止在高并发情况下多个请求同时查询数据库导致的缓存击穿问题，如果获取成功则返回用户对象和nil表示没有错误，并记录二次检查命中的日志
		err := s.cache.GetJSON(key, &u)
		if err == nil {
			logger.Info("double check hit cache")

			return &u, nil
		}
		logger.Info("double check cache miss")

		//二次检查未命中缓存，继续查询数据库获取用户个人信息
		userPtr, err := s.repo.GetUserByName("test")

		//如果查询数据库失败或者用户不存在则返回错误信息，并记录查询数据库失败的日志
		if err != nil || userPtr == nil {
			logger.Error("query DB failed for user")
			return nil, err
		}

		exp := time.Hour + time.Duration(rand.Intn(300))*time.Second //设置一个随机的过期时间，避免缓存雪崩问题，过期时间在1小时到1小时5分钟之间

		_ = s.cache.SetJSON(key, userPtr, exp) //将从数据库中获取的用户个人信息以JSON格式存储到Redis缓存中，设置一个随机的过期时间

		return userPtr, nil
	})

	//如果singleflight执行过程中发生错误则返回错误信息，并记录获取用户个人信息失败的日志
	if err != nil {
		logger.Error("singleflight get user profile failed")
		return nil, err
	}

	//记录获取用户个人信息成功的日志，包含用户ID信息
	logger.Info("get user profile success")
	return v.(*model.User), nil
}
