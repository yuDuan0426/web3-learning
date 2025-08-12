package sqlx

import "gorm.io/gorm"

type employee struct {
	ID         uint `gorm:"primaryKey"`
	Name       string
	Department string
	Salary     float64
}

type result struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

// 查询employees表中部门为“技术部”的员工信息，并映射到自定义的Employee结构切片中
func getEmployeeByTec(db *gorm.DB) *result {
	var employees []employee
	err := db.Table("employees").Where("department = ?", "技术部").Find(&employees).Error
	if err != nil {
		return nil // 处理错误
	}

	// 假设我们只需要第一个员工的信息
	if len(employees) > 0 {
		return &result{
			ID:   employees[0].ID,
			Name: employees[0].Name,
		}
	}
	return nil // 如果没有找到符合条件的员工
}

// 使用sqlx查询employees中工资最高的员工信息，将结果映射到Employee结构体中
func getHighestSalaryEmployee(db *gorm.DB) (*employee, error) {
	var emp employee
	err := db.Table("employees").Order("salary DESC").First(&emp).Error
	if err != nil {
		return nil, err // 返回错误
	}
	return &emp, nil // 返回找到的员工信息
}
