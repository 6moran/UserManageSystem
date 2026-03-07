package repo_mysql

import "GoWebUser/models/model"

type UserRepository interface {
	CreateUser(user *model.User) error
	GetByEmail(user *model.User) (*model.User, error)
	UpdateLastTime(user *model.User) error
	GetByUserID(user *model.User) (*model.User, error)
	GetLimitUsers(page, size int, status, keyword string) ([]*model.User, int, error)
	DeleteUserByID(user *model.User) error
	UpdateUserByID(user *model.User) error
}
