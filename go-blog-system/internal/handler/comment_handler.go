package handler

import (
	"net/http"
	"strconv"

	"github.com/duanyu/go-blog-system/internal/model"
	"github.com/duanyu/go-blog-system/internal/service"
	"github.com/gin-gonic/gin"
)

// CommentHandler 评论处理器
type CommentHandler struct {
	commentService service.CommentService
}

// NewCommentHandler 创建评论处理器
func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// Create 创建评论
func (h *CommentHandler) Create(c *gin.Context) {
	userID := GetUserIDFromContext(c)

	var req model.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := h.commentService.Create(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// Get 获取评论
func (h *CommentHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	comment, err := h.commentService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// Update 更新评论
func (h *CommentHandler) Update(c *gin.Context) {
	userID := GetUserIDFromContext(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	var req model.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := h.commentService.Update(id, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// Delete 删除评论
func (h *CommentHandler) Delete(c *gin.Context) {
	userID := GetUserIDFromContext(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	if err := h.commentService.Delete(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment deleted successfully"})
}

// GetByPost 获取文章的所有评论
func (h *CommentHandler) GetByPost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	comments, err := h.commentService.GetByPostID(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// RegisterRoutes 注册路由
func (h *CommentHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/posts/comments/:post_id", h.GetByPost)
	router.GET("/comments/:id", h.Get)

	authRouter := router.Group("/")
	authRouter.Use(AuthMiddleware())
	{
		authRouter.POST("/comments", h.Create)
		authRouter.PUT("/comments/:id", h.Update)
		authRouter.DELETE("/comments/:id", h.Delete)
	}
}