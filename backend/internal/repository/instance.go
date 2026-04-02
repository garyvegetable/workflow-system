package repository

import (
	"workflow-system/internal/domain/instance"
	"workflow-system/internal/domain/task"

	"gorm.io/gorm"
)

type InstanceRepository struct {
	db *gorm.DB
}

func NewInstanceRepository(db *gorm.DB) *InstanceRepository {
	return &InstanceRepository{db: db}
}

func (r *InstanceRepository) Create(inst *instance.WorkflowInstance) error {
	return r.db.Create(inst).Error
}

func (r *InstanceRepository) GetByID(id int64) (*instance.WorkflowInstance, error) {
	var inst instance.WorkflowInstance
	err := r.db.First(&inst, id).Error
	if err != nil {
		return nil, err
	}
	return &inst, nil
}

func (r *InstanceRepository) List(companyID int64) ([]instance.WorkflowInstance, error) {
	var instances []instance.WorkflowInstance
	err := r.db.Where("company_id = ?", companyID).Find(&instances).Error
	return instances, err
}

func (r *InstanceRepository) ListByInitiator(initiatorID int64) ([]instance.WorkflowInstance, error) {
	var instances []instance.WorkflowInstance
	err := r.db.Where("initiator_id = ?", initiatorID).Find(&instances).Error
	return instances, err
}

func (r *InstanceRepository) Update(inst *instance.WorkflowInstance) error {
	return r.db.Save(inst).Error
}

func (r *InstanceRepository) Cancel(id int64) error {
	return r.db.Model(&instance.WorkflowInstance{}).Where("id = ?", id).Update("status", 4).Error
}

func (r *InstanceRepository) CreateTask(t *task.ApprovalTask) error {
	return r.db.Create(t).Error
}

func (r *InstanceRepository) GetTaskByID(id int64) (*task.ApprovalTask, error) {
	var t task.ApprovalTask
	err := r.db.First(&t, id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *InstanceRepository) GetPendingTasks(assigneeID int64, companyID int64) ([]task.ApprovalTask, error) {
	var tasks []task.ApprovalTask
	err := r.db.Table("approval_task at").
		Select("at.*").
		Joins("LEFT JOIN workflow_instance wi ON wi.id = at.instance_id").
		Where("at.assignee_id = ? AND at.status = ? AND wi.company_id = ?", assigneeID, 1, companyID).
		Find(&tasks).Error
	return tasks, err
}

func (r *InstanceRepository) GetHandledTasks(assigneeID int64, companyID int64) ([]task.ApprovalTask, error) {
	var tasks []task.ApprovalTask
	err := r.db.Table("approval_task at").
		Select("at.*").
		Joins("LEFT JOIN workflow_instance wi ON wi.id = at.instance_id").
		Where("at.assignee_id = ? AND at.status != ? AND wi.company_id = ?", assigneeID, 1, companyID).
		Find(&tasks).Error
	return tasks, err
}

func (r *InstanceRepository) GetTasksByInstanceID(instanceID int64) ([]task.ApprovalTask, error) {
	var tasks []task.ApprovalTask
	err := r.db.Where("instance_id = ?", instanceID).Order("created_at").Find(&tasks).Error
	return tasks, err
}

func (r *InstanceRepository) UpdateTask(t *task.ApprovalTask) error {
	return r.db.Save(t).Error
}

func (r *InstanceRepository) UpdateTaskWithOptimisticLock(t *task.ApprovalTask, expectedVersion int) error {
	result := r.db.Model(&task.ApprovalTask{}).
		Where("id = ? AND version = ?", t.ID, expectedVersion).
		Updates(map[string]interface{}{
			"status":      t.Status,
			"action":      t.Action,
			"comment":     t.Comment,
			"completed_at": t.CompletedAt,
			"version":     expectedVersion + 1,
		})

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// DeleteTask deletes a task by ID
func (r *InstanceRepository) DeleteTask(taskID int64) error {
	return r.db.Delete(&task.ApprovalTask{}, taskID).Error
}

// GetOverduePendingTasks 获取已超时但仍未审批的任务
func (r *InstanceRepository) GetOverduePendingTasks(nowUnix int64) ([]task.ApprovalTask, error) {
	var tasks []task.ApprovalTask
	err := r.db.Where("status = ? AND deadline_at IS NOT NULL AND deadline_at < ?", 1, nowUnix).
		Find(&tasks).Error
	return tasks, err
}
