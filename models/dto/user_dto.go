package dto

import "time"

// UserDto 用户响应
type UserDto struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Status     int       `json:"status"`
	Role       string    `json:"role"`
	Avatar     string    `json:"avatar"` //头像
	CreateTime time.Time `json:"create_time"`
	LastTime   time.Time `json:"last_time"`
}

// AuthRequest 通用请求体
type AuthRequest struct {
	Email    string `json:"email" validator:"required,email,max=254"`
	Password string `json:"password" validator:"required,min=6,max=50"`
}

// AuthResponse 通用响应体
type AuthResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
