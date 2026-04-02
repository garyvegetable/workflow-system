package instance

import "encoding/json"

type WorkflowInstance struct {
	ID           int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	DefinitionID int64           `gorm:"not null" json:"definition_id"`
	CompanyID    int64           `gorm:"not null" json:"company_id"`
	Title        string          `gorm:"type:varchar(500);not null" json:"title"`
	InitiatorID  int64           `gorm:"not null" json:"initiator_id"`
	FormData     json.RawMessage `gorm:"type:jsonb;default:'{}'" json:"form_data"`
	Status       int16           `gorm:"default:0" json:"status"` // 0:草稿 1:审批中 2:已通过 3:已驳回 4:已撤回
	GraphData    json.RawMessage `gorm:"type:jsonb;default:'{}'" json:"graph_data"` // 保存流程定义快照
	CurrentNodes json.RawMessage `gorm:"type:jsonb;default:'[]'" json:"current_nodes"`
	StartedAt    *json.RawMessage `json:"started_at"`
	FinishedAt   *json.RawMessage `json:"finished_at"`
}

func (WorkflowInstance) TableName() string {
	return "workflow_instance"
}
