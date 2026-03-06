package services

import (
	"GoWebUser/models/dto"
)

type UserService interface {
	RegisterUser(reg dto.AuthRequest) error
	LoginUser(log dto.AuthRequest) (string, error)
	GetUserStatusByID(id int) (int, error)
	GetUsersByLimit(page, size int, status, keyword string) ([]*dto.UserDto, int, error)
	DeleteUserByID(id int) error
}
