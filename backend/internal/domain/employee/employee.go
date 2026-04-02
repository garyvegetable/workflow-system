package employee

type Employee struct {
	ID           int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyID    int64  `gorm:"not null" json:"company_id"`
	Username     string `gorm:"type:varchar(100);not null" json:"username"`
	Name         string `gorm:"type:varchar(100)" json:"name"`
	Email        string `gorm:"type:varchar(255)" json:"email"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	Level        string `gorm:"type:varchar(50)" json:"level"`
	Status       int16  `gorm:"default:1" json:"status"`
}

func (Employee) TableName() string {
	return "employee"
}

// EmployeeDepartment is the junction table for employee-department many-to-many relationship
type EmployeeDepartment struct {
	EmployeeID   int64 `gorm:"primaryKey"`
	DepartmentID int64 `gorm:"primaryKey"`
}

func (EmployeeDepartment) TableName() string {
	return "employee_department"
}
