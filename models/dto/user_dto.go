package dto

import "time"

// UserDto 用户响应
type UserDto struct {
	ID         int        `json:"id"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	Status     int        `json:"status"`
	Role       string     `json:"role"`
	Avatar     string     `json:"avatar"` //头像
	CreateTime *time.Time `json:"create_time"`
	LastTime   *time.Time `json:"last_time"`
}

// RegisterRequest 注册请求体
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=254"`
	Password string `json:"password" validate:"required,min=6,max=50"`
}

// LoginRequest 登录请求体
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse 通用响应体
type AuthResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// EditRequest 编辑请求体
type EditRequest struct {
	ID       int
	Username string `json:"username" validate:"required,max=50"`
	Password string `json:"password" validate:"omitempty,min=6,max=50"`
	Status   int    `json:"status" validate:"oneof=0 1"`
}
