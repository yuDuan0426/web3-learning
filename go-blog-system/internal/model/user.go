package model

import "time"

// User 用户模型
type User struct {
	ID        int       `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"` // 不在JSON中返回密码
	Avatar    *string   `db:"avatar" json:"avatar"`
	Bio       *string   `db:"bio" json:"bio"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// UserResponse 用户响应模型（不包含敏感信息）
type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    *string   `json:"avatar"`
	Bio       *string   `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse 转换为响应模型
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Avatar:    u.Avatar,
		Bio:       u.Bio,
		CreatedAt: u.CreatedAt,
	}
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string  `json:"username" binding:"required,min=3,max=50"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Avatar   *string `json:"avatar"`
	Bio      *string `json:"bio"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Email  *string `json:"email" binding:"omitempty,email"`
	Avatar *string `json:"avatar"`
	Bio    *string `json:"bio"`
}