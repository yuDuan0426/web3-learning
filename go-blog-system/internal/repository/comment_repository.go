package repository

import (
	"fmt"

	"github.com/duanyu/go-blog-system/internal/model"
	"github.com/jmoiron/sqlx"
)

// CommentRepository 评论仓库接口
type CommentRepository interface {
	Create(comment *model.Comment) error
	GetByID(id int) (*model.Comment, error)
	Update(comment *model.Comment) error
	Delete(id int) error
	GetByPostID(postID int) ([]model.Comment, error)
	GetReplies(commentID int) ([]model.Comment, error)
}

// commentRepository 评论仓库实现
type commentRepository struct {
	db *sqlx.DB
}

// NewCommentRepository 创建评论仓库
func NewCommentRepository(db *sqlx.DB) CommentRepository {
	return &commentRepository{db: db}
}

// Create 创建评论
func (r *commentRepository) Create(comment *model.Comment) error {
	query := `INSERT INTO comments (content, user_id, post_id, parent_id) 
			VALUES (?, ?, ?, ?)`

	result, err := r.db.Exec(query, comment.Content, comment.UserID, comment.PostID, comment.ParentID)
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	comment.ID = int(id)
	return nil
}

// GetByID 根据ID获取评论
func (r *commentRepository) GetByID(id int) (*model.Comment, error) {
	var comment model.Comment
	query := `SELECT * FROM comments WHERE id = ?`

	err := r.db.Get(&comment, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment by id: %w", err)
	}

	return &comment, nil
}

// Update 更新评论
func (r *commentRepository) Update(comment *model.Comment) error {
	query := `UPDATE comments SET content = ? WHERE id = ?`

	_, err := r.db.Exec(query, comment.Content, comment.ID)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

// Delete 删除评论
func (r *commentRepository) Delete(id int) error {
	query := `DELETE FROM comments WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}

// GetByPostID 获取文章的所有顶级评论（不包括回复）
func (r *commentRepository) GetByPostID(postID int) ([]model.Comment, error) {
	query := `SELECT * FROM comments WHERE post_id = ? AND parent_id IS NULL ORDER BY created_at DESC`

	var comments []model.Comment
	err := r.db.Select(&comments, query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by post id: %w", err)
	}

	return comments, nil
}

// GetReplies 获取评论的所有回复
func (r *commentRepository) GetReplies(commentID int) ([]model.Comment, error) {
	query := `SELECT * FROM comments WHERE parent_id = ? ORDER BY created_at ASC`

	var replies []model.Comment
	err := r.db.Select(&replies, query, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get replies: %w", err)
	}

	return replies, nil
}