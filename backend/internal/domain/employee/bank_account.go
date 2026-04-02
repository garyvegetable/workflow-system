package employee

type EmployeeBankAccount struct {
	ID           int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	EmployeeID   int64  `gorm:"not null" json:"employee_id"`
	BankName     string `gorm:"type:varchar(200);not null" json:"bank_name"`
	BankBranch   string `gorm:"type:varchar(200)" json:"bank_branch"`
	BankAccount  string `gorm:"type:varchar(100);not null" json:"bank_account"`
	AccountHolder string `gorm:"type:varchar(100);not null" json:"account_holder"`
	IsDefault    bool   `gorm:"default:false" json:"is_default"`
}

func (EmployeeBankAccount) TableName() string {
	return "employee_bank_account"
}
