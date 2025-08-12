package podcast

import "gorm.io/gorm"

// 一个用户可以有多篇文章，一片文章可以有多个评论
// 用户表
type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey;autoIncrement:true"`
	Username string `gorm:"uniqueIndex;not null;Size:100"`
	Password string `gorm:"not null;Size:100"`

	//关联文章  一对多    一个用户可以发表多篇文章
	Posts []Post `gorm:"foreignKey:AuthorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	//关联评论  一对多    一篇文章多条评论
	Comments []Comment `gorm:"foreignKey:PostID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// 文章表
type Post struct {
	gorm.Model
	ID uint `gorm:"primaryKey;autoIncrement:true"`
	// AuthorID 关联用户ID    多对一
	AuthorID uint   `gorm:"not null"`
	Title    string `gorm:"not null;Size:200"`
	Content  string `gorm:"not null;Size:5000"`

	// 关联用户，多对一，一个用户发表多篇文章
	Author User `gorm:"foreignKey:AuthorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	//关联评论，一对多，一篇文章可以有多条评论
	Comments []Comment `gorm:"foreignKey:PostID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// Comment表
type Comment struct {
	gorm.Model
	ID      uint   `gorm:"primaryKey;autoIncrement:true"`
	PostID  uint   `gorm:"not null;comment: '外键，关联post.id'"`
	UserID  uint   `gorm:"not null;comment: '外键，关联user.id'"`
	Content string `gorm:"not null;Size:500"`

	//关联文章，多对一    多条评论对应一篇文章
	Post Post `gorm:"foreignKey:PostID;comment: '关联的文章'"`
	//关联用户，多对一    多条评论对应一个用户
	User User `gorm:"foreignKey:UserID;comment: '关联的用户'"`
}

// 帮我插入一些数据，并关联起来
func InsertSampleData(db *gorm.DB) error {
	// 创建用户
	user1 := User{Username: "user1", Password: "pass1"}
	user2 := User{Username: "user2", Password: "pass2"}

	if err := db.Create(&user1).Error; err != nil {
		return err
	}
	if err := db.Create(&user2).Error; err != nil {
		return err
	}

	// 创建文章
	post1 := Post{AuthorID: user1.ID, Title: "Post by User 1", Content: "Content of post by user 1"}
	post2 := Post{AuthorID: user2.ID, Title: "Post by User 2", Content: "Content of post by user 2"}

	if err := db.Create(&post1).Error; err != nil {
		return err
	}
	if err := db.Create(&post2).Error; err != nil {
		return err
	}

	// 创建评论
	comment1 := Comment{PostID: post1.ID, UserID: user1.ID, Content: "Comment on post by user 1"}
	comment2 := Comment{PostID: post2.ID, UserID: user2.ID, Content: "Comment on post by user 2"}

	if err := db.Create(&comment1).Error; err != nil {
		return err
	}
	if err := db.Create(&comment2).Error; err != nil {
		return err
	}

	return nil
}
