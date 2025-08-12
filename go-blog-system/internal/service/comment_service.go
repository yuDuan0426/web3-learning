package service

import (
	"errors"
	"fmt"

	"github.com/duanyu/go-blog-system/internal/model"
	"github.com/duanyu/go-blog-system/internal/repository"
)

// CommentService 评论服务接口
type CommentService interface {
	Create(userID int, req *model.CreateCommentRequest) (*model.CommentResponse, error)
	GetByID(id int) (*model.CommentResponse, error)
	Update(id, userID int, req *model.UpdateCommentRequest) (*model.CommentResponse, error)
	Delete(id, userID int) error
	GetByPostID(postID int) ([]model.CommentResponse, error)
}

// commentService 评论服务实现
type commentService struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
	userRepo    repository.UserRepository
}

// NewCommentService 创建评论服务
func NewCommentService(
	commentRepo repository.CommentRepository,
	postRepo repository.PostRepository,
	userRepo repository.UserRepository,
) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		userRepo:    userRepo,
	}
}

// Create 创建评论
func (s *commentService) Create(userID int, req *model.CreateCommentRequest) (*model.CommentResponse, error) {
	// 检查文章是否存在
	_, err := s.postRepo.GetByID(req.PostID)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}

	// 如果有父评论，检查父评论是否存在
	if req.ParentID != nil {
		parentComment, err := s.commentRepo.GetByID(*req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent comment not found: %w", err)
		}

		// 确保父评论属于同一篇文章
		if parentComment.PostID != req.PostID {
			return nil, errors.New("parent comment does not belong to the specified post")
		}
	}

	// 创建评论
	comment := &model.Comment{
		Content:  req.Content,
		UserID:   userID,
		PostID:   req.PostID,
		ParentID: req.ParentID,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 构建响应
	response := &model.CommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		PostID:    comment.PostID,
		ParentID:  comment.ParentID,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		User: &model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Avatar:    user.Avatar,
			Bio:       user.Bio,
			CreatedAt: user.CreatedAt,
		},
	}

	return response, nil
}

// GetByID 根据ID获取评论
func (s *commentService) GetByID(id int) (*model.CommentResponse, error) {
	// 获取评论
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	// 获取用户信息
	user, err := s.userRepo.GetByID(comment.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 获取回复
	replies, err := s.commentRepo.GetReplies(comment.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get replies: %w", err)
	}

	// 构建响应
	response := &model.CommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		PostID:    comment.PostID,
		ParentID:  comment.ParentID,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		User: &model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Avatar:    user.Avatar,
			Bio:       user.Bio,
			CreatedAt: user.CreatedAt,
		},
	}

	// 添加回复
	if len(replies) > 0 {
		response.Replies = make([]model.Comment, len(replies))
		for i, reply := range replies {
			response.Replies[i] = reply
		}
	}

	return response, nil
}

// Update 更新评论
func (s *commentService) Update(id, userID int, req *model.UpdateCommentRequest) (*model.CommentResponse, error) {
	// 获取评论
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	// 检查权限
	if comment.UserID != userID {
		return nil, errors.New("you don't have permission to update this comment")
	}

	// 更新评论
	comment.Content = req.Content

	if err := s.commentRepo.Update(comment); err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	// 获取用户信息
	user, err := s.userRepo.GetByID(comment.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 构建响应
	response := &model.CommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		PostID:    comment.PostID,
		ParentID:  comment.ParentID,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		User: &model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Avatar:    user.Avatar,
			Bio:       user.Bio,
			CreatedAt: user.CreatedAt,
		},
	}

	return response, nil
}

// Delete 删除评论
func (s *commentService) Delete(id, userID int) error {
	// 获取评论
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get comment: %w", err)
	}

	// 检查权限
	if comment.UserID != userID {
		return errors.New("you don't have permission to delete this comment")
	}

	// 删除评论
	return s.commentRepo.Delete(id)
}

// GetByPostID 获取文章的所有评论
func (s *commentService) GetByPostID(postID int) ([]model.CommentResponse, error) {
	// 获取文章的所有顶级评论
	comments, err := s.commentRepo.GetByPostID(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	// 构建响应
	responses := make([]model.CommentResponse, len(comments))
	for i, comment := range comments {
		// 获取用户信息
		user, err := s.userRepo.GetByID(comment.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}

		// 获取回复
		replies, err := s.commentRepo.GetReplies(comment.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get replies: %w", err)
		}

		responses[i] = model.CommentResponse{
			ID:        comment.ID,
			Content:   comment.Content,
			PostID:    comment.PostID,
			ParentID:  comment.ParentID,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			User: &model.UserResponse{
				ID:        user.ID,
				Username:  user.Username,
				Email:     user.Email,
				Avatar:    user.Avatar,
				Bio:       user.Bio,
				CreatedAt: user.CreatedAt,
			},
		}

		// 添加回复
		if len(replies) > 0 {
			responses[i].Replies = make([]model.Comment, len(replies))
			for j, reply := range replies {
				responses[i].Replies[j] = reply
			}
		}
	}

	return responses, nil
}