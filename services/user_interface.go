package services

import (
	"GoWebUser/models/dto"
	"mime/multipart"
)

type UserService interface {
	RegisterUser(reg dto.RegisterRequest) error
	LoginUser(log dto.LoginRequest) (string, error)
	GetUserStatusByID(id int) (int, error)
	GetUsersByLimit(page, size int, status, keyword string) ([]*dto.UserDto, int, error)
	DeleteUserByID(id int) error
	GetUserRoleAndAvatar(id int) (*dto.UserDto, error)
	UpdateUserByID(req dto.EditRequest, file multipart.File, fileHeader *multipart.FileHeader) error
}
