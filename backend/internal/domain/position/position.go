package position

type Position struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyID int64 `gorm:"not null" json:"company_id"`
	Name     string `gorm:"type:varchar(50);not null" json:"name"`
	Code     string `gorm:"type:varchar(20);not null" json:"code"`
	Status   int16  `gorm:"default:1" json:"status"`
}

func (Position) TableName() string {
	return "position"
}
