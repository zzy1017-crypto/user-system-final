package repository

import "user-system-final/internal/model"

// MemoryRepo 是一个内存中的用户仓库实现，使用一个map来存储用户数据，提供GetUserByName方法来根据用户名获取用户信息
type MemoryRepo struct {
	users map[string]*model.User
}

// NewMemoryRepo 创建一个新的 MemoryRepo 实例，初始化一个包含测试用户的用户数据，并返回 MemoryRepo 对象
func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		users: map[string]*model.User{
			"test": {
				ID:       1,
				Username: "test",
				Password: "123456",
			},
		},
	}
}

// GetUserByName 根据用户名获取用户信息，如果用户存在则返回用户对象，否则返回nil，提供给服务层调用以实现用户认证和授权功能
func (r *MemoryRepo) GetUserByName(name string) (*model.User, error) {
	user, ok := r.users[name]
	if !ok {
		return nil, nil
	}
	return user, nil
}
