package book

import "gorm.io/gorm"

type book struct {
	ID     uint `gorm:"primaryKey"`
	Title  string
	Author string
	Price  float64
}

// 批量插入10条数据，5条的书籍价格大于50，其余5条小于50
func InsertBooks(db *gorm.DB) error {
	books := []book{
		{Title: "Book A", Author: "Author A", Price: 60.0},
		{Title: "Book B", Author: "Author B", Price: 70.0},
		{Title: "Book C", Author: "Author C", Price: 80.0},
		{Title: "Book D", Author: "Author D", Price: 90.0},
		{Title: "Book E", Author: "Author E", Price: 100.0},
		{Title: "Book F", Author: "Author F", Price: 40.0},
		{Title: "Book G", Author: "Author G", Price: 30.0},
		{Title: "Book H", Author: "Author H", Price: 20.0},
		{Title: "Book I", Author: "Author I", Price: 10.0},
		{Title: "Book J", Author: "Author J", Price: 5.0},
	}
	tx := db.Begin()
	for _, b := range books {
		if err := tx.Create(&b).Error; err != nil {
			tx.Rollback() // 如果插入失败，回滚事务
			return err
		}
	}
	return tx.Commit().Error // 提交事务
}

// 使用sqlx执行一个复杂查询，查询价格大于50元的书籍，将结果映射到Book结构体切片中，确保类型安全
func getBooksByPrice(db *gorm.DB) ([]book, error) {
	var books []book
	err := db.Where("price > ?", 50).Find(&books).Error
	if err != nil {
		return nil, err // 返回错误
	}
	return books, nil // 返回找到的书籍信息
}
