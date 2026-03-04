package services

import (
	"GoWebUser/models/dto"
	"GoWebUser/models/model"
	"GoWebUser/repositories/repo_mysql"
	"GoWebUser/utils"
	"errors"
	"log"
	"strings"
)

var (
	ErrorsEmailExists    = errors.New("email already exists")
	ErrorsEmailNotExists = errors.New("email not exists")
	ErrorsWrongPassword  = errors.New("wrong password")
	ErrorsIdNotExists    = errors.New("id not exists")
)

type UserServiceImpl struct {
	UserRepo repo_mysql.UserRepository
}

func NewUserService(ur repo_mysql.UserRepository) UserService {
	return &UserServiceImpl{UserRepo: ur}
}

// 处理注册业务
func (u *UserServiceImpl) RegisterUser(r dto.AuthRequest) error {
	//转换为model结构
	reg := &model.User{
		Username: "用户" + strings.Split(r.Email, "@")[0],
		Email:    r.Email,
		Password: r.Password,
	}

	//调用数据库来创建用户
	err := u.UserRepo.CreateUser(reg)
	if err != nil {
		if errors.Is(err, repo_mysql.ErrEmailExists) {
			return ErrorsEmailExists
		}
		return err
	}
	return nil
}

// 处理登录业务
func (u *UserServiceImpl) LoginUser(l dto.AuthRequest) (string, error) {
	//转换为model结构
	login := &model.User{
		Email:    l.Email,
		Password: l.Password,
	}

	//调用数据库函数
	user, err := u.UserRepo.GetByEmail(login)
	if err != nil {
		if errors.Is(err, repo_mysql.ErrEmailNotExists) {
			return "", ErrorsEmailNotExists
		}
		return "", err
	}
	if user.Password != l.Password {
		return "", ErrorsWrongPassword
	}

	//更新最后一次登录的时间
	err = u.UserRepo.UpdateLastTime(user)
	if err != nil {
		//这里即使失败登录也成功了
		log.Printf("UpdateLastTime failed,err:%v", err)
	}

	return utils.GetToken(user.ID)
}

// GetUserStatusByID 封装在服务层，用于中间件调用
func (u *UserServiceImpl) GetUserStatusByID(id int) (int, error) {
	user, err := u.UserRepo.GetByUserID(&model.User{ID: id})
	if err != nil {
		if errors.Is(err, repo_mysql.ErrIdNotExists) {
			return 0, ErrorsIdNotExists
		}
		return 0, err
	}
	return user.Status, nil
}

// GetUsersByLimit 处理查用户业务
func (u *UserServiceImpl) GetUsersByLimit(page, size int, status, keyword string) ([]*dto.UserDto, int, error) {
	userList, num, err := u.UserRepo.GetLimitUsers(page, size, status, keyword)
	if err != nil {
		return nil, 0, err
	}
	var users []*dto.UserDto
	for _, user := range userList {
		ud := &dto.UserDto{
			Username:   user.Username,
			Email:      user.Email,
			Status:     user.Status,
			Role:       user.Role,
			Avatar:     user.Avatar,
			CreateTime: user.CreateTime,
			LastTime:   user.LastTime,
		}
		users = append(users, ud)
	}
	return users, num, nil
}
