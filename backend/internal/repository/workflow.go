package repository

import (
	"workflow-system/internal/domain/workflow"

	"gorm.io/gorm"
)

type WorkflowRepository struct {
	db *gorm.DB
}

func NewWorkflowRepository(db *gorm.DB) *WorkflowRepository {
	return &WorkflowRepository{db: db}
}

func (r *WorkflowRepository) Create(wf *workflow.WorkflowDefinition) error {
	return r.db.Create(wf).Error
}

func (r *WorkflowRepository) GetByID(id int64) (*workflow.WorkflowDefinition, error) {
	var wf workflow.WorkflowDefinition
	err := r.db.First(&wf, id).Error
	if err != nil {
		return nil, err
	}
	return &wf, nil
}

func (r *WorkflowRepository) List(companyID int64) ([]workflow.WorkflowDefinition, error) {
	var definitions []workflow.WorkflowDefinition
	err := r.db.Where("company_id = ?", companyID).Find(&definitions).Error
	return definitions, err
}

func (r *WorkflowRepository) Update(wf *workflow.WorkflowDefinition) error {
	updates := map[string]interface{}{}
	if wf.Name != "" {
		updates["name"] = wf.Name
	}
	if wf.Code != "" {
		updates["code"] = wf.Code
	}
	if len(wf.GraphData) > 0 {
		updates["graph_data"] = wf.GraphData
	}
	if wf.FormFields != nil {
		updates["form_fields"] = wf.FormFields
	}
	if wf.Version != 0 {
		updates["version"] = wf.Version
	}
	if wf.Status != 0 {
		updates["status"] = wf.Status
	}
	return r.db.Model(&workflow.WorkflowDefinition{}).Where("id = ?", wf.ID).Updates(updates).Error
}

func (r *WorkflowRepository) Delete(id int64) error {
	return r.db.Delete(&workflow.WorkflowDefinition{}, id).Error
}

func (r *WorkflowRepository) Publish(id int64) error {
	return r.db.Model(&workflow.WorkflowDefinition{}).Where("id = ?", id).Update("status", 2).Error
}

func (r *WorkflowRepository) Disable(id int64) error {
	return r.db.Model(&workflow.WorkflowDefinition{}).Where("id = ?", id).Update("status", 3).Error
}

func (r *WorkflowRepository) Enable(id int64) error {
	return r.db.Model(&workflow.WorkflowDefinition{}).Where("id = ?", id).Update("status", 2).Error
}

func (r *WorkflowRepository) CopyToCompany(id int64, targetCompanyID int64) (*workflow.WorkflowDefinition, error) {
	var source workflow.WorkflowDefinition
	if err := r.db.First(&source, id).Error; err != nil {
		return nil, err
	}

	newWF := workflow.WorkflowDefinition{
		CompanyID:  targetCompanyID,
		Code:       source.Code + "_copy",
		Name:       source.Name + " (副本)",
		Version:    1,
		GraphData:  source.GraphData,
		FormFields: source.FormFields,
		Status:     1,
	}

	if err := r.db.Create(&newWF).Error; err != nil {
		return nil, err
	}
	return &newWF, nil
}
