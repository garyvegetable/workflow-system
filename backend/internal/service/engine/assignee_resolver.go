package engine

import (
	"workflow-system/internal/repository"
)

type AssigneeResolver struct {
	deptRepo   *repository.DepartmentRepository
	employeeRepo *repository.EmployeeRepository
}

func NewAssigneeResolver(deptRepo *repository.DepartmentRepository, employeeRepo *repository.EmployeeRepository) *AssigneeResolver {
	return &AssigneeResolver{
		deptRepo:   deptRepo,
		employeeRepo: employeeRepo,
	}
}

func (r *AssigneeResolver) Resolve(nodeType string, nodeData map[string]interface{}, context map[string]interface{}) (int64, error) {
	switch nodeType {
	case "user":
		// 指定用户
		if userID, ok := nodeData["assigneeId"].(float64); ok {
			return int64(userID), nil
		}

	case "department_leader":
		// 部门负责人
		if deptID, ok := nodeData["departmentId"].(float64); ok {
			dept, err := r.deptRepo.GetByID(int64(deptID))
			if err != nil {
				return 0, err
			}
			if dept.LeaderID != nil {
				return *dept.LeaderID, nil
			}
		}

	case "role":
		// 角色级别
		// 实际实现需要根据级别找到对应的审批人
	}

	return 0, nil
}

func (r *AssigneeResolver) GetNextApprover(departmentID int64, currentStep int) (int64, error) {
	chain, err := r.deptRepo.GetApprovalChain(departmentID)
	if err != nil {
		return 0, err
	}

	for _, step := range chain {
		if step.StepOrder > currentStep {
			return step.EmployeeID, nil
		}
	}

	return 0, nil // 没有更多审批人
}
