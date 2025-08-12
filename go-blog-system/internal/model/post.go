package model

import "time"

// PostStatus 文章状态
type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusArchived  PostStatus = "archived"
)

// Post 文章模型
type Post struct {
	ID        int        `db:"id" json:"id"`
	Title     string     `db:"title" json:"title"`
	Content   string     `db:"content" json:"content"`
	UserID    int        `db:"user_id" json:"user_id"`
	Status    PostStatus `db:"status" json:"status"`
	ViewCount int        `db:"view_count" json:"view_count"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	// 关联字段（不在数据库中）
	User      *UserResponse `db:"-" json:"user,omitempty"`
	Tags      []Tag         `db:"-" json:"tags,omitempty"`
	Comments  []Comment     `db:"-" json:"comments,omitempty"`
}

// PostResponse 文章响应模型
type PostResponse struct {
	ID        int           `json:"id"`
	Title     string        `json:"title"`
	Content   string        `json:"content"`
	Status    PostStatus    `json:"status"`
	ViewCount int           `json:"view_count"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	User      *UserResponse `json:"user,omitempty"`
	Tags      []Tag         `json:"tags,omitempty"`
}

// ToResponse 转换为响应模型
func (p *Post) ToResponse() PostResponse {
	return PostResponse{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		Status:    p.Status,
		ViewCount: p.ViewCount,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		User:      p.User,
		Tags:      p.Tags,
	}
}

// CreatePostRequest 创建文章请求
type CreatePostRequest struct {
	Title   string     `json:"title" binding:"required"`
	Content string     `json:"content" binding:"required"`
	Status  PostStatus `json:"status" binding:"omitempty"`
	TagIDs  []int      `json:"tag_ids" binding:"omitempty"`
}

// UpdatePostRequest 更新文章请求
type UpdatePostRequest struct {
	Title   *string     `json:"title" binding:"omitempty"`
	Content *string     `json:"content" binding:"omitempty"`
	Status  *PostStatus `json:"status" binding:"omitempty"`
	TagIDs  []int       `json:"tag_ids" binding:"omitempty"`
}

// PostQuery 文章查询参数
type PostQuery struct {
	UserID  *int       `form:"user_id"`
	Status  PostStatus `form:"status"`
	TagID   *int       `form:"tag_id"`
	Keyword string     `form:"keyword"`
	Page    int        `form:"page,default=1"`
	PerPage int        `form:"per_page,default=10"`
}