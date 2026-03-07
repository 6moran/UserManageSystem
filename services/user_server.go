package services

import (
	"GoWebUser/models/dto"
	"GoWebUser/models/model"
	"GoWebUser/repositories/repo_mysql"
	"GoWebUser/utils"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrorsEmailExists    = errors.New("email already exists")
	ErrorsEmailNotExists = errors.New("email not exists")
	ErrorsWrongPassword  = errors.New("wrong password")
	ErrorsIdNotExists    = errors.New("id not exists")
	ErrorsStatusNotMatch = errors.New("the status does not match")
	ErrorsJustImages     = errors.New("just need images")
	ErrorsNotMoreThan2MB = errors.New("images must less than 2MB ")
)

type UserServiceImpl struct {
	UserRepo repo_mysql.UserRepository
}

func NewUserService(ur repo_mysql.UserRepository) UserService {
	return &UserServiceImpl{UserRepo: ur}
}

// RegisterUser 处理注册业务
func (u *UserServiceImpl) RegisterUser(r dto.RegisterRequest) error {
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

// LoginUser 处理登录业务
func (u *UserServiceImpl) LoginUser(l dto.LoginRequest) (string, error) {
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
	if user.Status != 1 {
		return "", ErrorsStatusNotMatch
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

// GetUsersByLimit 处理分页查全部用户业务
func (u *UserServiceImpl) GetUsersByLimit(page, size int, status, keyword string) ([]*dto.UserDto, int, error) {
	userList, num, err := u.UserRepo.GetLimitUsers(page, size, status, keyword)
	if err != nil {
		return nil, 0, err
	}
	var users []*dto.UserDto
	for _, user := range userList {
		ud := &dto.UserDto{
			ID:         user.ID,
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

// DeleteUserByID 处理删除用户业务
func (u *UserServiceImpl) DeleteUserByID(id int) error {
	err := u.UserRepo.DeleteUserByID(&model.User{ID: id})
	if err != nil {
		return err
	}
	return nil
}

// GetUserRoleAndAvatar 处理查询当前用户信息业务
func (u *UserServiceImpl) GetUserRoleAndAvatar(id int) (*dto.UserDto, error) {
	user, err := u.UserRepo.GetByUserID(&model.User{ID: id})
	if err != nil {
		return nil, err
	}
	return &dto.UserDto{
		ID:     id,
		Email:  user.Email,
		Role:   user.Role,
		Avatar: user.Avatar,
	}, nil
}

// UpdateUserByID 处理根据ID修改用户信息业务
func (u *UserServiceImpl) UpdateUserByID(req dto.EditRequest, file multipart.File, fileHeader *multipart.FileHeader) error {
	avatar := ""
	// 头像文件处理
	if file != nil && fileHeader != nil {
		// 类型限制
		if !strings.HasPrefix(fileHeader.Header.Get("Content-Type"), "image/") {
			return ErrorsJustImages
		}
		// 大小限制
		if fileHeader.Size > 2*1024*1024 {
			return ErrorsNotMoreThan2MB
		}

		// 读取文件内容生成 MD5
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		hash := fmt.Sprintf("%x", md5.Sum(fileBytes))
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if ext == "" {
			ext = ".png" // 默认扩展名
		}

		// 创建目录
		os.MkdirAll("static/img", 0755)

		filename := hash + ext
		filePath := filepath.Join("static/img", filename)

		// 文件不存在才保存
		if _, err = os.Stat(filePath); os.IsNotExist(err) {
			err = os.WriteFile(filePath, fileBytes, 0644)
			if err != nil {
				return err
			}
		}

		avatar = "/static/img/" + filename
	}

	err := u.UserRepo.UpdateUserByID(&model.User{
		ID:       req.ID,
		Username: req.Username,
		Password: req.Password,
		Status:   req.Status,
		Avatar:   avatar,
	})
	if err != nil {
		return err
	}
	return nil
}
