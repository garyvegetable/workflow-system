package service

import (
	"encoding/json"
	"fmt"
	"workflow-system/internal/domain/instance"
	"workflow-system/internal/domain/workflow"
	"workflow-system/internal/repository"
)

type WorkflowService struct {
	repo        *repository.WorkflowRepository
	instanceRepo *repository.InstanceRepository
}

func NewWorkflowService(repo *repository.WorkflowRepository, instanceRepo *repository.InstanceRepository) *WorkflowService {
	return &WorkflowService{repo: repo, instanceRepo: instanceRepo}
}

func (s *WorkflowService) Create(wf *workflow.WorkflowDefinition) error {
	return s.repo.Create(wf)
}

func (s *WorkflowService) GetByID(id int64) (*workflow.WorkflowDefinition, error) {
	return s.repo.GetByID(id)
}

func (s *WorkflowService) List(companyID int64) ([]workflow.WorkflowDefinition, error) {
	return s.repo.List(companyID)
}

func (s *WorkflowService) Update(wf *workflow.WorkflowDefinition) error {
	return s.repo.Update(wf)
}

func (s *WorkflowService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *WorkflowService) Publish(id int64) error {
	return s.repo.Publish(id)
}

func (s *WorkflowService) Copy(id int64, targetCompanyID int64) (*workflow.WorkflowDefinition, error) {
	return s.repo.CopyToCompany(id, targetCompanyID)
}

func (s *WorkflowService) Disable(id int64) error {
	return s.repo.Disable(id)
}

func (s *WorkflowService) CreateInstance(inst *instance.WorkflowInstance, formData map[string]interface{}) error {
	inst.FormData, _ = json.Marshal(formData)
	inst.Status = 1 // 审批中
	return s.instanceRepo.Create(inst)
}

func (s *WorkflowService) GetInstance(id int64) (*instance.WorkflowInstance, error) {
	return s.instanceRepo.GetByID(id)
}

func (s *WorkflowService) CancelInstance(id int64, userID int64) error {
	// 获取实例并验证是否是申请人本人
	inst, err := s.instanceRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("流程实例不存在")
	}
	if inst.InitiatorID != userID {
		return fmt.Errorf("只有申请人可以撤回")
	}
	if inst.Status != 1 {
		return fmt.Errorf("只有审批中的申请可以撤回")
	}
	return s.instanceRepo.Cancel(id)
}

func (s *WorkflowService) GetMyApplications(initiatorID int64) ([]instance.WorkflowInstance, error) {
	return s.instanceRepo.ListByInitiator(initiatorID)
}
