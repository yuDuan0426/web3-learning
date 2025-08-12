package repository

import (
	"fmt"
	"strings"

	"github.com/duanyu/go-blog-system/internal/model"
	"github.com/jmoiron/sqlx"
)

// PostRepository 文章仓库接口
type PostRepository interface {
	Create(post *model.Post) error
	GetByID(id int) (*model.Post, error)
	Update(post *model.Post) error
	Delete(id int) error
	List(query *model.PostQuery) ([]model.Post, error)
	Count(query *model.PostQuery) (int, error)
	IncrementViewCount(id int) error
	AddTags(postID int, tagIDs []int) error
	RemoveTags(postID int) error
	GetPostTags(postID int) ([]model.Tag, error)
}

// postRepository 文章仓库实现
type postRepository struct {
	db *sqlx.DB
}

// NewPostRepository 创建文章仓库
func NewPostRepository(db *sqlx.DB) PostRepository {
	return &postRepository{db: db}
}

// Create 创建文章
func (r *postRepository) Create(post *model.Post) error {
	query := `INSERT INTO posts (title, content, user_id, status) 
			VALUES (?, ?, ?, ?)`

	result, err := r.db.Exec(query, post.Title, post.Content, post.UserID, post.Status)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	post.ID = int(id)
	return nil
}

// GetByID 根据ID获取文章
func (r *postRepository) GetByID(id int) (*model.Post, error) {
	var post model.Post
	query := `SELECT * FROM posts WHERE id = ?`

	err := r.db.Get(&post, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post by id: %w", err)
	}

	return &post, nil
}

// Update 更新文章
func (r *postRepository) Update(post *model.Post) error {
	query := `UPDATE posts SET title = ?, content = ?, status = ? WHERE id = ?`

	_, err := r.db.Exec(query, post.Title, post.Content, post.Status, post.ID)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	return nil
}

// Delete 删除文章
func (r *postRepository) Delete(id int) error {
	query := `DELETE FROM posts WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

// List 获取文章列表
func (r *postRepository) List(query *model.PostQuery) ([]model.Post, error) {
	baseQuery := `SELECT * FROM posts WHERE 1=1`
	whereClause, args := r.buildWhereClause(query)

	offset := (query.Page - 1) * query.PerPage
	finalQuery := baseQuery + whereClause + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`

	args = append(args, query.PerPage, offset)

	var posts []model.Post
	err := r.db.Select(&posts, finalQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}

	return posts, nil
}

// Count 获取文章总数
func (r *postRepository) Count(query *model.PostQuery) (int, error) {
	baseQuery := `SELECT COUNT(*) FROM posts WHERE 1=1`
	whereClause, args := r.buildWhereClause(query)

	finalQuery := baseQuery + whereClause

	var count int
	err := r.db.Get(&count, finalQuery, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to count posts: %w", err)
	}

	return count, nil
}

// IncrementViewCount 增加文章浏览量
func (r *postRepository) IncrementViewCount(id int) error {
	query := `UPDATE posts SET view_count = view_count + 1 WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to increment view count: %w", err)
	}

	return nil
}

// AddTags 添加文章标签
func (r *postRepository) AddTags(postID int, tagIDs []int) error {
	if len(tagIDs) == 0 {
		return nil
	}

	// 构建批量插入语句
	values := make([]string, 0, len(tagIDs))
	args := make([]interface{}, 0, len(tagIDs)*2)

	for _, tagID := range tagIDs {
		values = append(values, "(?, ?)")
		args = append(args, postID, tagID)
	}

	query := fmt.Sprintf(
		"INSERT INTO post_tags (post_id, tag_id) VALUES %s",
		strings.Join(values, ", "),
	)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to add tags to post: %w", err)
	}

	return nil
}

// RemoveTags 移除文章所有标签
func (r *postRepository) RemoveTags(postID int) error {
	query := `DELETE FROM post_tags WHERE post_id = ?`

	_, err := r.db.Exec(query, postID)
	if err != nil {
		return fmt.Errorf("failed to remove tags from post: %w", err)
	}

	return nil
}

// GetPostTags 获取文章标签
func (r *postRepository) GetPostTags(postID int) ([]model.Tag, error) {
	query := `
		SELECT t.* 
		FROM tags t
		JOIN post_tags pt ON t.id = pt.tag_id
		WHERE pt.post_id = ?
	`

	var tags []model.Tag
	err := r.db.Select(&tags, query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post tags: %w", err)
	}

	return tags, nil
}

// buildWhereClause 构建WHERE子句
func (r *postRepository) buildWhereClause(query *model.PostQuery) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if query.UserID != nil {
		conditions = append(conditions, "user_id = ?")
		args = append(args, *query.UserID)
	}

	if query.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, query.Status)
	}

	if query.TagID != nil {
		conditions = append(conditions, "id IN (SELECT post_id FROM post_tags WHERE tag_id = ?)")
		args = append(args, *query.TagID)
	}

	if query.Keyword != "" {
		conditions = append(conditions, "(title LIKE ? OR content LIKE ?)")
		keyword := "%" + query.Keyword + "%"
		args = append(args, keyword, keyword)
	}

	if len(conditions) == 0 {
		return "", args
	}

	return " AND " + strings.Join(conditions, " AND "), args
}