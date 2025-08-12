# Go博客系统

一个基于Go语言开发的博客系统，包含用户管理、文章管理、评论管理等功能。

## 项目结构

项目采用领域驱动设计和清晰分层架构：

```
.
├── cmd/                # 应用程序入口
│   └── api/            # API服务入口
├── config/             # 配置文件
├── docs/               # 文档
├── internal/           # 内部包
│   ├── handler/        # HTTP处理器
│   ├── model/          # 数据模型
│   ├── repository/     # 数据仓库
│   └── service/        # 业务服务
├── migrations/         # 数据库迁移文件
├── pkg/                # 公共包
│   ├── database/       # 数据库连接
│   └── logger/         # 日志工具
└── scripts/            # 脚本工具
```

## 功能特性

- 用户管理：注册、登录、个人资料管理
- 文章管理：创建、编辑、删除、查看文章
- 评论系统：发表评论、回复评论
- 标签管理：创建标签、为文章添加标签

## 技术栈

- Go语言
- Gin Web框架
- MySQL数据库
- SQLX库
- JWT认证
- Viper配置管理
- Logrus日志库

## 快速开始

### 前置条件

- Go 1.16+
- MySQL 5.7+

### 安装

1. 克隆项目

```bash
git clone https://github.com/yourusername/go-blog-system.git
cd go-blog-system
```

2. 安装依赖

```bash
go mod tidy
```

3. 配置数据库

编辑 `config/config.yaml` 文件，设置数据库连接信息：

```yaml
database:
  driver: "mysql"
  host: "localhost"
  port: 3306
  username: "your_username"
  password: "your_password"
  dbname: "go_blog"
  params: "charset=utf8mb4&parseTime=True&loc=Local"
```

4. 创建数据库

```sql
CREATE DATABASE go_blog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

5. 运行数据库迁移

```bash
go run scripts/migrate.go
```

6. 启动应用

```bash
go run cmd/api/main.go
```

应用将在 http://localhost:8080 上运行。

## API文档

### 用户相关

- `POST /api/register` - 注册新用户
- `POST /api/login` - 用户登录
- `GET /api/profile` - 获取当前用户信息
- `PUT /api/profile` - 更新用户信息
- `GET /api/users/:id` - 获取指定用户信息
- `GET /api/users` - 获取用户列表

### 文章相关

- `POST /api/posts` - 创建文章
- `GET /api/posts/:id` - 获取文章详情
- `PUT /api/posts/:id` - 更新文章
- `DELETE /api/posts/:id` - 删除文章
- `GET /api/posts` - 获取文章列表

### 评论相关

- `POST /api/comments` - 创建评论
- `GET /api/comments/:id` - 获取评论详情
- `PUT /api/comments/:id` - 更新评论
- `DELETE /api/comments/:id` - 删除评论
- `GET /api/posts/:post_id/comments` - 获取文章的所有评论

### 标签相关

- `POST /api/tags` - 创建标签
- `GET /api/tags/:id` - 获取标签详情
- `GET /api/tags` - 获取所有标签
- `DELETE /api/tags/:id` - 删除标签

## 许可证

本项目采用 MIT 许可证。