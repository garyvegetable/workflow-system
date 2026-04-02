package service

import (
	"fmt"
	"workflow-system/internal/domain/department"
	"workflow-system/internal/repository"
)

type DepartmentService struct {
	repo          *repository.DepartmentRepository
	employeeRepo  *repository.EmployeeRepository
}

func NewDepartmentService(repo *repository.DepartmentRepository, employeeRepo *repository.EmployeeRepository) *DepartmentService {
	return &DepartmentService{repo: repo, employeeRepo: employeeRepo}
}

func (s *DepartmentService) Create(dept *department.Department) error {
	return s.repo.Create(dept)
}

func (s *DepartmentService) GetByID(id int64) (*department.Department, error) {
	return s.repo.GetByID(id)
}

func (s *DepartmentService) List(companyID int64) ([]department.Department, error) {
	return s.repo.List(companyID)
}

func (s *DepartmentService) Count(companyID int64) (int64, error) {
	return s.repo.Count(companyID)
}

func (s *DepartmentService) GetAllChildIDs(deptID int64) []int64 {
	return s.repo.GetAllChildIDs(deptID)
}

func (s *DepartmentService) GetTree(companyID int64) ([]*department.TreeNode, error) {
	return s.repo.GetTree(companyID)
}

func (s *DepartmentService) Update(dept *department.Department) error {
	return s.repo.Update(dept)
}

// UpdateFields 更新指定字段（支持 null 值）
func (s *DepartmentService) UpdateFields(id int64, fields map[string]interface{}) error {
	return s.repo.UpdateFields(id, fields)
}

// Delete 删除部门
// transferToDeptID: 员工转移目标部门，0 表示不转移
func (s *DepartmentService) Delete(id int64, transferToDeptID int64) error {
	return s.repo.Delete(id, transferToDeptID)
}

func (s *DepartmentService) GetApprovalChain(deptID int64) ([]department.DepartmentApprovalChain, error) {
	return s.repo.GetApprovalChain(deptID)
}

// ValidateApprovalChainEmployees 验证审批链中的员工是否属于该部门
func (s *DepartmentService) ValidateApprovalChainEmployees(deptID int64, chain []department.DepartmentApprovalChain) error {
	if len(chain) == 0 {
		return nil
	}

	// 获取部门的所有子部门ID（员工可能在子部门中）
	childDeptIDs := s.repo.GetAllChildIDs(deptID)

	for i, step := range chain {
		if step.EmployeeID == 0 {
			continue
		}
		// 获取员工所属的部门列表
		empDepts, err := s.employeeRepo.GetDepartments(step.EmployeeID)
		if err != nil {
			return fmt.Errorf("无法获取员工 %d 的部门信息: %w", step.EmployeeID, err)
		}

		// 检查员工是否属于该部门或其子部门
		belongs := false
		for _, empDeptID := range empDepts {
			for _, childDeptID := range childDeptIDs {
				if empDeptID == childDeptID {
					belongs = true
					break
				}
			}
			if belongs {
				break
			}
		}

		if !belongs {
			return fmt.Errorf("第 %d 步审批人（员工ID: %d）不属于该部门或子部门", i+1, step.EmployeeID)
		}
	}
	return nil
}

// SetApprovalChain 设置审批链（先验证再设置）
func (s *DepartmentService) SetApprovalChain(deptID int64, chain []department.DepartmentApprovalChain) error {
	// 先验证
	if err := s.ValidateApprovalChainEmployees(deptID, chain); err != nil {
		return err
	}
	return s.repo.SetApprovalChain(deptID, chain)
}
