package supplier

type Supplier struct {
	ID         int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyID  int64  `gorm:"not null" json:"company_id"`
	Code       string `gorm:"type:varchar(50);not null" json:"code"`
	Name       string `gorm:"type:varchar(200);not null" json:"name"`
	Contact    string `gorm:"type:varchar(100)" json:"contact"`
	Phone      string `gorm:"type:varchar(50)" json:"phone"`
	Email      string `gorm:"type:varchar(255)" json:"email"`
	Address    string `gorm:"type:varchar(500)" json:"address"`
	BankName   string `gorm:"type:varchar(200)" json:"bank_name"`
	BankAccount string `gorm:"type:varchar(100)" json:"bank_account"`
	TaxNumber  string `gorm:"type:varchar(50)" json:"tax_number"`
	Status     int16  `gorm:"default:1" json:"status"`
}

func (Supplier) TableName() string {
	return "supplier"
}
