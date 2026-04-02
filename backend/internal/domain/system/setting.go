package system

// SystemSetting 键值对形式的系统设置
type SystemSetting struct {
	Key   string `gorm:"primaryKey;size:100" json:"key"`
	Value string `gorm:"type:text" json:"value"`
}

func (SystemSetting) TableName() string {
	return "system_setting"
}

// 预定义的设置键
const (
	KeySMTPHost     = "smtp_host"
	KeySMTPPort     = "smtp_port"
	KeySMTPUser     = "smtp_user"
	KeySMTPPassword = "smtp_password"
	KeySMTPFrom     = "smtp_from"
)
