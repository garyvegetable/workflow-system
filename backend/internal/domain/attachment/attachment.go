package attachment

type Attachment struct {
	ID         int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	InstanceID int64  `gorm:"not null" json:"instance_id"`
	FieldName  string `gorm:"type:varchar(100)" json:"field_name"`
	FileName   string `gorm:"type:varchar(255);not null" json:"file_name"`
	FilePath   string `gorm:"type:varchar(500);not null" json:"file_path"`
	FileSize   int64  `json:"file_size"`
	MimeType   string `gorm:"type:varchar(100)" json:"mime_type"`
}

func (Attachment) TableName() string {
	return "attachment"
}
