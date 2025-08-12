package main

import (
	"fmt"
	"log"

	"github.com/duanyu/go-blog-system/internal/handler"
	"github.com/duanyu/go-blog-system/internal/repository"
	"github.com/duanyu/go-blog-system/internal/service"
	"github.com/duanyu/go-blog-system/pkg/database"
	"github.com/duanyu/go-blog-system/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// 加载配置
	if err := loadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	if err := logger.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// 连接数据库
	db, err := database.InitFromViper()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 设置Gin模式
	mode := viper.GetString("app.mode")
	if mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 创建仓库
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	tagRepo := repository.NewTagRepository(db)

	// 创建服务
	userService := service.NewUserService(userRepo)
	postService := service.NewPostService(postRepo, userRepo, tagRepo)
	commentService := service.NewCommentService(commentRepo, postRepo, userRepo)
	tagService := service.NewTagService(tagRepo)

	// 创建处理器
	userHandler := handler.NewUserHandler(userService)
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)
	tagHandler := handler.NewTagHandler(tagService)

	// 注册路由
	api := r.Group("/api")
	{
		userHandler.RegisterRoutes(api)
		postHandler.RegisterRoutes(api)
		commentHandler.RegisterRoutes(api)
		tagHandler.RegisterRoutes(api)
	}

	// 启动服务器
	port := viper.GetInt("app.port")
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// loadConfig 加载配置
func loadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	return viper.ReadInConfig()
}