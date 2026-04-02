package task

type ApprovalTask struct {
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	InstanceID  int64  `gorm:"not null" json:"instance_id"`
	NodeID      string `gorm:"type:varchar(50);not null" json:"node_id"`
	NodeName    string `gorm:"type:varchar(100);not null" json:"node_name"`
	AssigneeID  int64  `gorm:"not null" json:"assignee_id"`
	Status      int16  `gorm:"default:1" json:"status"` // 1:待审批 2:已审批 3:已驳回
	Action      string `gorm:"type:varchar(20)" json:"action"` // approve/reject
	Comment     string `gorm:"type:text" json:"comment"`
	DeadlineAt  *int64 `gorm:"default:null" json:"deadline_at"` // 超时时间戳
	CompletedAt *int64 `json:"completed_at"`
	Version     int    `gorm:"default:0" json:"version"` // 乐观锁
	IsOverdue   bool   `gorm:"-" json:"is_overdue"` // 计算字段，不存储
}

func (ApprovalTask) TableName() string {
	return "approval_task"
}
