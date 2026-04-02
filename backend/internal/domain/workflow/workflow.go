package workflow

import "encoding/json"

type WorkflowDefinition struct {
	ID          int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyID   int64           `gorm:"not null" json:"company_id"`
	Code        string          `gorm:"type:varchar(50);not null" json:"code"`
	Name        string          `gorm:"type:varchar(200);not null" json:"name"`
	Version     int             `gorm:"default:1" json:"version"`
	GraphData   json.RawMessage `gorm:"type:jsonb;default:'{}'" json:"graph_data"`
	FormFields  json.RawMessage `gorm:"type:jsonb;default:'[]'" json:"form_fields"`
	Status      int16           `gorm:"default:1" json:"status"` // 1:草稿 2:已发布 3:禁用
}

func (WorkflowDefinition) TableName() string {
	return "workflow_definition"
}
