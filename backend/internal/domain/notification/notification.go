package notification

type Notification struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64  `gorm:"not null" json:"user_id"`
	Title     string `gorm:"type:varchar(200);not null" json:"title"`
	Content   string `gorm:"type:text" json:"content"`
	Type      string `gorm:"type:varchar(50)" json:"type"`
	IsRead    bool   `gorm:"default:false" json:"is_read"`
}

func (Notification) TableName() string {
	return "notification"
}
