package repository

import "user-system-final/internal/model"

type MemoryRepo struct {
	users map[string]*model.User
}

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

func (r *MemoryRepo) GetUserByName(name string) (*model.User, error) {
	user, ok := r.users[name]
	if !ok {
		return nil, nil
	}
	return user, nil
}
