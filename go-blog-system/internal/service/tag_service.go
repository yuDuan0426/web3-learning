package service

import (
	"errors"
	"fmt"

	"github.com/duanyu/go-blog-system/internal/model"
	"github.com/duanyu/go-blog-system/internal/repository"
)

// TagService 标签服务接口
type TagService interface {
	Create(req *model.CreateTagRequest) (*model.Tag, error)
	GetByID(id int) (*model.Tag, error)
	List() ([]model.Tag, error)
	Delete(id int) error
}

// tagService 标签服务实现
type tagService struct {
	tagRepo repository.TagRepository
}

// NewTagService 创建标签服务
func NewTagService(tagRepo repository.TagRepository) TagService {
	return &tagService{tagRepo: tagRepo}
}

// Create 创建标签
func (s *tagService) Create(req *model.CreateTagRequest) (*model.Tag, error) {
	// 检查标签名是否已存在
	_, err := s.tagRepo.GetByName(req.Name)
	if err == nil {
		return nil, errors.New("tag name already exists")
	}

	// 创建标签
	tag := &model.Tag{
		Name: req.Name,
	}

	if err := s.tagRepo.Create(tag); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return tag, nil
}

// GetByID 根据ID获取标签
func (s *tagService) GetByID(id int) (*model.Tag, error) {
	return s.tagRepo.GetByID(id)
}

// List 获取所有标签
func (s *tagService) List() ([]model.Tag, error) {
	return s.tagRepo.List()
}

// Delete 删除标签
func (s *tagService) Delete(id int) error {
	return s.tagRepo.Delete(id)
}