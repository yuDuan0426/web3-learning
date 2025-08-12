package main

import (
	"fmt"

	"github.com/we3-learning/task3/podcast"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:Dy950426@tcp(localhost:3306)/podcast?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic("failed to connect database")
	}
	//lesson.Run(db)

	//db.AutoMigrate(&podcast.User{})
	//db.AutoMigrate(&podcast.Post{})
	//db.AutoMigrate(&podcast.Comment{})
	//
	//podcast.InsertSampleData(db)

	//dsn := "root:Dy950426@tcp(localhost:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := sqlx.Connect("mysql", dsn)
	//if err != nil {
	//	panic("failed to connect database")
	//}
	//defer db.Close()
	//
	//err = db.Ping()
	//if err != nil {
	//	panic("failed to ping database")
	//}
	//fmt.Println("Successfully connected to database")

	//查找用户发布的所有文章以及对应的评论信息
	var users []podcast.User
	err = db.Preload("Posts.Comments").Where("username = ?", "user1").Find(&users).Error
	if err != nil {
		panic(err)
	}
	fmt.Println(users)

	//使用Gorm查询评论数量最多的文章信息。
	//select * from comments where poist_id  = (SELECT count(post_id) as comments_count FROM comments GROUP BY post_id ORDER BY count(post_id) DESC LIMIT 1);

	var postWithMostComments podcast.Post
	err = db.Model(&podcast.Comment{}).
		Select("post_id, COUNT(*) as comments_count").
		Joins("JOIN posts ON posts.id = comments.post_id").
		Group("post_id").
		Order("comments_count DESC").
		Limit(1).
		Scan(&postWithMostComments).Error
	if err != nil {
		panic(err)
	}
	fmt.Println(postWithMostComments)

}

// 创建钩子函数,要求在
