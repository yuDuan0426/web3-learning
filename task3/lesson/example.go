package lesson

import (
	"fmt"

	"gorm.io/gorm"
)

type Student struct {
	ID    uint `gorm:"primaryKey";autoIncrement:true`
	Name  string
	Age   int
	Grade string
}

//创建两个表 account和transaction
// account表包含字段ID（主键，自增）、balance（余额）
// transaction表包含字段ID（主键）、from_account_id(转出id)，to_account_id(转入id)，amount(转账金额)

type account struct {
	ID      uint `gorm:"primaryKey"`
	Balance float64
}

type transaction struct {
	ID            uint `gorm:"primaryKey"`
	FromAccountID uint
	ToAccountID   uint
	Amount        float64
}

func Run(db *gorm.DB) {
	// Migrate the schema
	db.AutoMigrate(&Student{})

	//向student表中插入一条新纪录，学生姓名为“张三”，年龄为20，年级为“三年级”
	student := Student{Name: "张三", Age: 20, Grade: "三年级"}
	result := db.Create(&student)
	fmt.Println("error msg:", result.Error)
	fmt.Println("execute result:", result.RowsAffected)

	////查询students表中年龄大于18岁的学生信息
	db.Select(&Student{}).Where("age > ?", 18).Find(&[]Student{})

	//更新students表中姓名为“张三”的学生的年级为四年级
	db.Select(&Student{}).Where("name =?", "张三").Update("grade", "四年级")

	//删除student表中年龄小于15岁的学生记录
	db.Where("age<?", 15).Delete(&Student{})

	//编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
	tx := db.Begin()
	accountA := account{ID: 1, Balance: 500} // 假设账户 A 的 ID 为 1，余额为 500
	accountB := account{ID: 2, Balance: 300} // 假设账户 B 的 ID 为 2，余额为 300
	amount := 100.0
	if accountA.Balance >= amount {
		// 扣除账户 A 的余额
		accountA.Balance -= amount
		if err := tx.Save(&accountA).Error; err != nil {
			tx.Rollback()
			fmt.Println("Error updating account A:", err)
			return
		}

		// 增加账户 B 的余额
		accountB.Balance += amount
		if err := tx.Save(&accountB).Error; err != nil {
			tx.Rollback()
			fmt.Println("Error updating account B:", err)
			return
		}
	}

}
