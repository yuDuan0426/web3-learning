package repository

import (
	"fmt"

	"github.com/duanyu/go-blog-system/internal/model"
	"github.com/jmoiron/sqlx"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	Create(user *model.User) error
	GetByID(id int) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	Delete(id int) error
	List(page, perPage int) ([]model.User, error)
	Count() (int, error)
}

// userRepository 用户仓库实现
type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(user *model.User) error {
	query := `INSERT INTO users (username, email, password, avatar, bio) 
			VALUES (?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.Avatar, user.Bio)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = int(id)
	return nil
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(id int) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE id = ?`

	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE username = ?`

	err := r.db.Get(&user, query, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE email = ?`

	err := r.db.Get(&user, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// Update 更新用户
func (r *userRepository) Update(user *model.User) error {
	query := `UPDATE users SET email = ?, password = ?, avatar = ?, bio = ? WHERE id = ?`

	_, err := r.db.Exec(query, user.Email, user.Password, user.Avatar, user.Bio, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete 删除用户
func (r *userRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// List 获取用户列表
func (r *userRepository) List(page, perPage int) ([]model.User, error) {
	offset := (page - 1) * perPage
	query := `SELECT * FROM users LIMIT ? OFFSET ?`

	var users []model.User
	err := r.db.Select(&users, query, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// Count 获取用户总数
func (r *userRepository) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users`

	err := r.db.Get(&count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}