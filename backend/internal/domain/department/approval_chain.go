package department

type DepartmentApprovalChain struct {
	ID           int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	DepartmentID int64 `gorm:"not null" json:"department_id"`
	EmployeeID   int64 `gorm:"not null" json:"employee_id"`
	StepOrder    int   `gorm:"not null" json:"step_order"`
}

func (DepartmentApprovalChain) TableName() string {
	return "department_approval_chain"
}
