package handler

import (
	"net/http"
	"strconv"

	"github.com/duanyu/go-blog-system/internal/model"
	"github.com/duanyu/go-blog-system/internal/service"
	"github.com/gin-gonic/gin"
)

// PostHandler 文章处理器
type PostHandler struct {
	postService service.PostService
}

// NewPostHandler 创建文章处理器
func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

// Create 创建文章
func (h *PostHandler) Create(c *gin.Context) {
	userID := GetUserIDFromContext(c)

	var req model.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := h.postService.Create(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// Get 获取文章
func (h *PostHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	// 增加浏览量
	if c.Query("view") == "true" {
		if err := h.postService.IncrementViewCount(id); err != nil {
			// 记录错误但不中断请求
			c.Error(err)
		}
	}

	post, err := h.postService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// Update 更新文章
func (h *PostHandler) Update(c *gin.Context) {
	userID := GetUserIDFromContext(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var req model.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := h.postService.Update(id, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

// Delete 删除文章
func (h *PostHandler) Delete(c *gin.Context) {
	userID := GetUserIDFromContext(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	if err := h.postService.Delete(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post deleted successfully"})
}

// List 获取文章列表
func (h *PostHandler) List(c *gin.Context) {
	var query model.PostQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认值
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PerPage <= 0 {
		query.PerPage = 10
	}

	posts, total, err := h.postService.List(&query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"meta": gin.H{
			"total":    total,
			"page":     query.Page,
			"per_page": query.PerPage,
		},
	})
}

// RegisterRoutes 注册路由
func (h *PostHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/posts", h.List)
	router.GET("/posts/:id", h.Get)

	authRouter := router.Group("/")
	authRouter.Use(AuthMiddleware())
	{
		authRouter.POST("/posts", h.Create)
		authRouter.PUT("/posts/:id", h.Update)
		authRouter.DELETE("/posts/:id", h.Delete)
	}
}