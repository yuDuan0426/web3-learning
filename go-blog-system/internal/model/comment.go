package model

import "time"

// Comment 评论模型
type Comment struct {
	ID        int        `db:"id" json:"id"`
	Content   string     `db:"content" json:"content"`
	UserID    int        `db:"user_id" json:"user_id"`
	PostID    int        `db:"post_id" json:"post_id"`
	ParentID  *int       `db:"parent_id" json:"parent_id"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	// 关联字段（不在数据库中）
	User      *UserResponse `db:"-" json:"user,omitempty"`
	Replies   []Comment     `db:"-" json:"replies,omitempty"`
}

// CommentResponse 评论响应模型
type CommentResponse struct {
	ID        int           `json:"id"`
	Content   string        `json:"content"`
	PostID    int           `json:"post_id"`
	ParentID  *int          `json:"parent_id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	User      *UserResponse `json:"user,omitempty"`
	Replies   []Comment     `json:"replies,omitempty"`
}

// ToResponse 转换为响应模型
func (c *Comment) ToResponse() CommentResponse {
	return CommentResponse{
		ID:        c.ID,
		Content:   c.Content,
		PostID:    c.PostID,
		ParentID:  c.ParentID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		User:      c.User,
		Replies:   c.Replies,
	}
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required"`
	PostID   int    `json:"post_id" binding:"required"`
	ParentID *int   `json:"parent_id" binding:"omitempty"`
}

// UpdateCommentRequest 更新评论请求
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}