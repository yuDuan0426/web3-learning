package repository

import (
	"fmt"

	"github.com/duanyu/go-blog-system/internal/model"
	"github.com/jmoiron/sqlx"
)

// TagRepository 标签仓库接口
type TagRepository interface {
	Create(tag *model.Tag) error
	GetByID(id int) (*model.Tag, error)
	GetByName(name string) (*model.Tag, error)
	List() ([]model.Tag, error)
	Delete(id int) error
}

// tagRepository 标签仓库实现
type tagRepository struct {
	db *sqlx.DB
}

// NewTagRepository 创建标签仓库
func NewTagRepository(db *sqlx.DB) TagRepository {
	return &tagRepository{db: db}
}

// Create 创建标签
func (r *tagRepository) Create(tag *model.Tag) error {
	query := `INSERT INTO tags (name) VALUES (?)`

	result, err := r.db.Exec(query, tag.Name)
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	tag.ID = int(id)
	return nil
}

// GetByID 根据ID获取标签
func (r *tagRepository) GetByID(id int) (*model.Tag, error) {
	var tag model.Tag
	query := `SELECT * FROM tags WHERE id = ?`

	err := r.db.Get(&tag, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag by id: %w", err)
	}

	return &tag, nil
}

// GetByName 根据名称获取标签
func (r *tagRepository) GetByName(name string) (*model.Tag, error) {
	var tag model.Tag
	query := `SELECT * FROM tags WHERE name = ?`

	err := r.db.Get(&tag, query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag by name: %w", err)
	}

	return &tag, nil
}

// List 获取所有标签
func (r *tagRepository) List() ([]model.Tag, error) {
	var tags []model.Tag
	query := `SELECT * FROM tags ORDER BY name ASC`

	err := r.db.Select(&tags, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	return tags, nil
}

// Delete 删除标签
func (r *tagRepository) Delete(id int) error {
	query := `DELETE FROM tags WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}