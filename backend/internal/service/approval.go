package service

import (
	"fmt"
	"workflow-system/internal/domain/task"
	"workflow-system/internal/repository"
	"workflow-system/internal/service/engine"
)

type ApprovalService struct {
	instanceRepo *repository.InstanceRepository
	workflowRepo *repository.WorkflowRepository
	engine       *engine.WorkflowEngine
}

func NewApprovalService(instanceRepo *repository.InstanceRepository, workflowRepo *repository.WorkflowRepository, wfEngine *engine.WorkflowEngine) *ApprovalService {
	return &ApprovalService{
		instanceRepo: instanceRepo,
		workflowRepo: workflowRepo,
		engine:       wfEngine,
	}
}

func (s *ApprovalService) ListPending(assigneeID int64, companyID int64) ([]task.ApprovalTask, error) {
	return s.instanceRepo.GetPendingTasks(assigneeID, companyID)
}

func (s *ApprovalService) ListHandled(assigneeID int64, companyID int64) ([]task.ApprovalTask, error) {
	return s.instanceRepo.GetHandledTasks(assigneeID, companyID)
}

func (s *ApprovalService) Approve(taskID int64, approverID int64, comment string) error {
	t, err := s.instanceRepo.GetTaskByID(taskID)
	if err != nil {
		return err
	}

	// 验证审批人权限
	if t.AssigneeID != approverID {
		return fmt.Errorf("您不是该任务的审批人，无权审批")
	}

	// 调用流程引擎处理审批结果，推进流程
	return s.engine.ProcessApproval(taskID, approverID, "approve", comment)
}

func (s *ApprovalService) Reject(taskID int64, approverID int64, comment string) error {
	t, err := s.instanceRepo.GetTaskByID(taskID)
	if err != nil {
		return err
	}

	// 验证审批人权限
	if t.AssigneeID != approverID {
		return fmt.Errorf("您不是该任务的审批人，无权审批")
	}

	// 调用流程引擎处理审批结果，推进流程
	return s.engine.ProcessApproval(taskID, approverID, "reject", comment)
}

func (s *ApprovalService) Transfer(taskID int64, newAssigneeID int64) error {
	t, err := s.instanceRepo.GetTaskByID(taskID)
	if err != nil {
		return err
	}

	t.AssigneeID = newAssigneeID
	return s.instanceRepo.UpdateTask(t)
}

func (s *ApprovalService) GetHistory(instanceID int64) ([]task.ApprovalTask, error) {
	return s.instanceRepo.GetTasksByInstanceID(instanceID)
}

// AddApprover 加签：为任务添加审批人
func (s *ApprovalService) AddApprover(taskID int64, newApproverID int64, approverID int64) error {
	t, err := s.instanceRepo.GetTaskByID(taskID)
	if err != nil {
		return err
	}

	// 只能为待审批状态的任务加签
	if t.Status != 1 {
		return fmt.Errorf("只能对待审批状态的任务加签")
	}

	// 创建新的审批任务
	newTask := &task.ApprovalTask{
		InstanceID: t.InstanceID,
		NodeID:    t.NodeID,
		NodeName:  t.NodeName,
		AssigneeID: newApproverID,
		Status:    1,
	}
	return s.instanceRepo.CreateTask(newTask)
}

// RemoveApprover 减签：移除待审批的审批人（只能移除未审批的任务）
func (s *ApprovalService) RemoveApprover(taskID int64, targetAssigneeID int64, approverID int64) error {
	t, err := s.instanceRepo.GetTaskByID(taskID)
	if err != nil {
		return err
	}

	// 只能减签待审批状态的任务
	if t.Status != 1 {
		return fmt.Errorf("只能对待审批状态的任务减签")
	}

	// 查找该节点上该审批人的所有待审批任务并删除
	tasks, err := s.instanceRepo.GetTasksByInstanceID(t.InstanceID)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if task.NodeID == t.NodeID && task.AssigneeID == targetAssigneeID && task.Status == 1 {
			// 删除该待审批任务
			s.instanceRepo.DeleteTask(task.ID)
		}
	}
	return nil
}
