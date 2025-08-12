package service

import (
	"errors"
	"fmt"

	"github.com/duanyu/go-blog-system/internal/model"
	"github.com/duanyu/go-blog-system/internal/repository"
)

// PostService 文章服务接口
type PostService interface {
	Create(userID int, req *model.CreatePostRequest) (*model.PostResponse, error)
	GetByID(id int) (*model.PostResponse, error)
	Update(id, userID int, req *model.UpdatePostRequest) (*model.PostResponse, error)
	Delete(id, userID int) error
	List(query *model.PostQuery) ([]model.PostResponse, int, error)
	IncrementViewCount(id int) error
}

// postService 文章服务实现
type postService struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
	tagRepo  repository.TagRepository
}

// NewPostService 创建文章服务
func NewPostService(
	postRepo repository.PostRepository,
	userRepo repository.UserRepository,
	tagRepo repository.TagRepository,
) PostService {
	return &postService{
		postRepo: postRepo,
		userRepo: userRepo,
		tagRepo:  tagRepo,
	}
}

// Create 创建文章
func (s *postService) Create(userID int, req *model.CreatePostRequest) (*model.PostResponse, error) {
	// 检查用户是否存在
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 创建文章
	status := req.Status
	if status == "" {
		status = model.PostStatusDraft
	}

	post := &model.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
		Status:  status,
	}

	if err := s.postRepo.Create(post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// 添加标签
	if len(req.TagIDs) > 0 {
		if err := s.postRepo.AddTags(post.ID, req.TagIDs); err != nil {
			return nil, fmt.Errorf("failed to add tags: %w", err)
		}
	}

	// 获取标签
	tags, err := s.postRepo.GetPostTags(post.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post tags: %w", err)
	}

	// 构建响应
	response := &model.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Status:    post.Status,
		ViewCount: post.ViewCount,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		User: &model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Avatar:    user.Avatar,
			Bio:       user.Bio,
			CreatedAt: user.CreatedAt,
		},
		Tags: tags,
	}

	return response, nil
}

// GetByID 根据ID获取文章
func (s *postService) GetByID(id int) (*model.PostResponse, error) {
	// 获取文章
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	// 获取作者
	user, err := s.userRepo.GetByID(post.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 获取标签
	tags, err := s.postRepo.GetPostTags(post.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post tags: %w", err)
	}

	// 构建响应
	response := &model.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Status:    post.Status,
		ViewCount: post.ViewCount,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		User: &model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Avatar:    user.Avatar,
			Bio:       user.Bio,
			CreatedAt: user.CreatedAt,
		},
		Tags: tags,
	}

	return response, nil
}

// Update 更新文章
func (s *postService) Update(id, userID int, req *model.UpdatePostRequest) (*model.PostResponse, error) {
	// 获取文章
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	// 检查权限
	if post.UserID != userID {
		return nil, errors.New("you don't have permission to update this post")
	}

	// 更新字段
	if req.Title != nil {
		post.Title = *req.Title
	}

	if req.Content != nil {
		post.Content = *req.Content
	}

	if req.Status != nil {
		post.Status = *req.Status
	}

	// 更新文章
	if err := s.postRepo.Update(post); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	// 更新标签
	if req.TagIDs != nil {
		// 先删除所有标签
		if err := s.postRepo.RemoveTags(post.ID); err != nil {
			return nil, fmt.Errorf("failed to remove tags: %w", err)
		}

		// 添加新标签
		if len(req.TagIDs) > 0 {
			if err := s.postRepo.AddTags(post.ID, req.TagIDs); err != nil {
				return nil, fmt.Errorf("failed to add tags: %w", err)
			}
		}
	}

	// 获取标签
	tags, err := s.postRepo.GetPostTags(post.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post tags: %w", err)
	}

	// 获取作者
	user, err := s.userRepo.GetByID(post.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 构建响应
	response := &model.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Status:    post.Status,
		ViewCount: post.ViewCount,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		User: &model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Avatar:    user.Avatar,
			Bio:       user.Bio,
			CreatedAt: user.CreatedAt,
		},
		Tags: tags,
	}

	return response, nil
}

// Delete 删除文章
func (s *postService) Delete(id, userID int) error {
	// 获取文章
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}

	// 检查权限
	if post.UserID != userID {
		return errors.New("you don't have permission to delete this post")
	}

	// 删除文章
	return s.postRepo.Delete(id)
}

// List 获取文章列表
func (s *postService) List(query *model.PostQuery) ([]model.PostResponse, int, error) {
	// 获取文章列表
	posts, err := s.postRepo.List(query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list posts: %w", err)
	}

	// 获取文章总数
	count, err := s.postRepo.Count(query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	// 构建响应
	responses := make([]model.PostResponse, len(posts))
	for i, post := range posts {
		// 获取作者
		user, err := s.userRepo.GetByID(post.UserID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get user: %w", err)
		}

		// 获取标签
		tags, err := s.postRepo.GetPostTags(post.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get post tags: %w", err)
		}

		responses[i] = model.PostResponse{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			Status:    post.Status,
			ViewCount: post.ViewCount,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
			User: &model.UserResponse{
				ID:        user.ID,
				Username:  user.Username,
				Email:     user.Email,
				Avatar:    user.Avatar,
				Bio:       user.Bio,
				CreatedAt: user.CreatedAt,
			},
			Tags: tags,
		}
	}

	return responses, count, nil
}

// IncrementViewCount 增加文章浏览量
func (s *postService) IncrementViewCount(id int) error {
	return s.postRepo.IncrementViewCount(id)
}