package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/duanyu/go-blog-system/internal/model"
	"github.com/duanyu/go-blog-system/internal/repository"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务接口
type UserService interface {
	Register(req *model.CreateUserRequest) (*model.UserResponse, error)
	Login(req *model.LoginRequest) (string, *model.UserResponse, error)
	GetByID(id int) (*model.UserResponse, error)
	Update(id int, req *model.UpdateUserRequest) (*model.UserResponse, error)
	Delete(id int) error
	List(page, perPage int) ([]model.UserResponse, int, error)
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// Register 注册用户
func (s *userService) Register(req *model.CreateUserRequest) (*model.UserResponse, error) {
	// 检查用户名是否已存在
	_, err := s.userRepo.GetByUsername(req.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	_, err = s.userRepo.GetByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Avatar:   req.Avatar,
		Bio:      req.Bio,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Bio:       user.Bio,
		CreatedAt: user.CreatedAt,
	}, nil
}

// Login 用户登录
func (s *userService) Login(req *model.LoginRequest) (string, *model.UserResponse, error) {
	// 获取用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return "", nil, errors.New("invalid username or password")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", nil, errors.New("invalid username or password")
	}

	// 生成JWT令牌
	token, err := generateJWT(user.ID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Bio:       user.Bio,
		CreatedAt: user.CreatedAt,
	}, nil
}

// GetByID 根据ID获取用户
func (s *userService) GetByID(id int) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Bio:       user.Bio,
		CreatedAt: user.CreatedAt,
	}, nil
}

// Update 更新用户
func (s *userService) Update(id int, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 更新字段
	if req.Email != nil {
		// 检查邮箱是否已被其他用户使用
		existingUser, err := s.userRepo.GetByEmail(*req.Email)
		if err == nil && existingUser.ID != id {
			return nil, errors.New("email already in use")
		}
		user.Email = *req.Email
	}

	if req.Avatar != nil {
		user.Avatar = req.Avatar
	}

	if req.Bio != nil {
		user.Bio = req.Bio
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Bio:       user.Bio,
		CreatedAt: user.CreatedAt,
	}, nil
}

// Delete 删除用户
func (s *userService) Delete(id int) error {
	return s.userRepo.Delete(id)
}

// List 获取用户列表
func (s *userService) List(page, perPage int) ([]model.UserResponse, int, error) {
	users, err := s.userRepo.List(page, perPage)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	count, err := s.userRepo.Count()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	userResponses := make([]model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Avatar:    user.Avatar,
			Bio:       user.Bio,
			CreatedAt: user.CreatedAt,
		}
	}

	return userResponses, count, nil
}

// generateJWT 生成JWT令牌
func generateJWT(userID int) (string, error) {
	jwtSecret := viper.GetString("app.jwt_secret")
	expiration := viper.GetInt("app.jwt_expiration")

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * time.Duration(expiration)).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}