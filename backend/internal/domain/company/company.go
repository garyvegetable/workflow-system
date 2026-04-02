package company

type Company struct {
	ID         int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Code       string `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name       string `gorm:"type:varchar(200);not null" json:"name"`
	ShortName  string `gorm:"type:varchar(100)" json:"short_name"`
	Status     int16  `gorm:"default:1" json:"status"` // 1:正常 2:禁用
	SchemaName string `gorm:"type:varchar(100)" json:"schema_name"`
}

func (Company) TableName() string {
	return "company"
}
