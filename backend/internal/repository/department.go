package repository

import (
	"workflow-system/internal/domain/department"
	"workflow-system/internal/domain/employee"

	"gorm.io/gorm"
)

type DepartmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) Create(dept *department.Department) error {
	return r.db.Create(dept).Error
}

func (r *DepartmentRepository) GetByID(id int64) (*department.Department, error) {
	var dept department.Department
	err := r.db.First(&dept, id).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *DepartmentRepository) List(companyID int64) ([]department.Department, error) {
	var departments []department.Department
	err := r.db.Where("company_id = ?", companyID).Find(&departments).Error
	return departments, err
}

func (r *DepartmentRepository) Count(companyID int64) (int64, error) {
	var count int64
	err := r.db.Model(&department.Department{}).Where("company_id = ?", companyID).Count(&count).Error
	return count, err
}

func (r *DepartmentRepository) GetTree(companyID int64) ([]*department.TreeNode, error) {
	var departments []department.Department
	err := r.db.Where("company_id = ?", companyID).Order("id").Find(&departments).Error
	if err != nil {
		return nil, err
	}

	// 构建树形结构
	nodeMap := make(map[int64]*department.TreeNode)
	var roots []*department.TreeNode

	// 先将所有部门转为节点
	for i := range departments {
		nodeMap[departments[i].ID] = &department.TreeNode{
			ID:        departments[i].ID,
			CompanyID: departments[i].CompanyID,
			ParentID:  departments[i].ParentID,
			Name:      departments[i].Name,
			LeaderID:  departments[i].LeaderID,
			Status:    departments[i].Status,
			Children:  []*department.TreeNode{},
		}
	}

	// 再构建父子关系
	for i := range departments {
		node := nodeMap[departments[i].ID]
		if departments[i].ParentID != nil && *departments[i].ParentID > 0 {
			if parent, ok := nodeMap[*departments[i].ParentID]; ok {
				parent.Children = append(parent.Children, node)
			} else {
				// 父节点不存在，当作根节点
				roots = append(roots, node)
			}
		} else {
			roots = append(roots, node)
		}
	}

	return roots, nil
}

func (r *DepartmentRepository) Update(dept *department.Department) error {
	// 只更新非零值字段，避免覆盖未提供的字段
	updates := map[string]interface{}{}
	if dept.Name != "" {
		updates["name"] = dept.Name
	}
	if dept.ParentID != nil {
		updates["parent_id"] = *dept.ParentID
	}
	if dept.LeaderID != nil {
		updates["leader_id"] = *dept.LeaderID
	}
	if dept.Status != 0 {
		updates["status"] = dept.Status
	}
	return r.db.Model(&department.Department{}).Where("id = ?", dept.ID).Updates(updates).Error
}

// UpdateFields 更新指定字段（支持显式 null 值）
func (r *DepartmentRepository) UpdateFields(id int64, fields map[string]interface{}) error {
	return r.db.Model(&department.Department{}).Where("id = ?", id).Updates(fields).Error
}

// Delete 删除部门及其子部门，员工转移到目标部门
// transferToDeptID: 员工转移目标部门，0 表示不转移直接删除
func (r *DepartmentRepository) Delete(id int64, transferToDeptID int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 获取所有子部门 ID（包括自己）
		allDeptIDs := r.getAllChildIDs(tx, id)

		// 2. 如果有员工需要转移
		if transferToDeptID > 0 {
			// 将这些部门的员工关系转移到目标部门
			// 只转移那些在要被删除的部门树中的员工
			if err := tx.Exec(`
				DELETE FROM employee_department
				WHERE department_id IN ? AND employee_id NOT IN (
					SELECT employee_id FROM employee_department WHERE department_id = ?
				)
			`, allDeptIDs, transferToDeptID).Error; err != nil {
				return err
			}
			// 添加目标部门的关系
			var empIDs []int64
			tx.Model(&employee.EmployeeDepartment{}).
				Where("department_id IN ?", allDeptIDs).
				Pluck("employee_id", &empIDs)
			for _, empID := range empIDs {
				// 检查是否已有该员工与目标部门的关系
				var exists int64
				tx.Model(&employee.EmployeeDepartment{}).
					Where("employee_id = ? AND department_id = ?", empID, transferToDeptID).
					Count(&exists)
				if exists == 0 {
					if err := tx.Create(&employee.EmployeeDepartment{
						EmployeeID:   empID,
						DepartmentID: transferToDeptID,
					}).Error; err != nil {
						return err
					}
				}
			}
		}

		// 3. 删除所有相关部门的审批链
		if err := tx.Where("department_id IN ?", allDeptIDs).Delete(&department.DepartmentApprovalChain{}).Error; err != nil {
			return err
		}

		// 4. 删除子部门（不包括自己，让外部处理）
		if len(allDeptIDs) > 1 {
			var childIDs []int64
			for _, deptID := range allDeptIDs {
				if deptID != id {
					childIDs = append(childIDs, deptID)
				}
			}
			if err := tx.Where("id IN ?", childIDs).Delete(&department.Department{}).Error; err != nil {
				return err
			}
		}

		// 5. 删除自己
		return tx.Delete(&department.Department{}, id).Error
	})
}

// getAllChildIDs 递归获取所有子部门 ID（包含自己，内部使用事务）
func (r *DepartmentRepository) getAllChildIDs(tx *gorm.DB, parentID int64) []int64 {
	var ids []int64
	ids = append(ids, parentID)

	var childIDs []int64
	tx.Model(&department.Department{}).Where("parent_id = ?", parentID).Pluck("id", &childIDs)
	for _, childID := range childIDs {
		ids = append(ids, r.getAllChildIDs(tx, childID)...)
	}
	return ids
}

// GetAllChildIDs 递归获取所有子部门 ID（包含自己）
func (r *DepartmentRepository) GetAllChildIDs(parentID int64) []int64 {
	var ids []int64
	ids = append(ids, parentID)

	var childIDs []int64
	r.db.Model(&department.Department{}).Where("parent_id = ?", parentID).Pluck("id", &childIDs)
	for _, childID := range childIDs {
		ids = append(ids, r.GetAllChildIDs(childID)...)
	}
	return ids
}

func (r *DepartmentRepository) GetApprovalChain(deptID int64) ([]department.DepartmentApprovalChain, error) {
	var chain []department.DepartmentApprovalChain
	err := r.db.Where("department_id = ?", deptID).Order("step_order").Find(&chain).Error
	return chain, err
}

func (r *DepartmentRepository) SetApprovalChain(deptID int64, chain []department.DepartmentApprovalChain) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除旧的
		if err := tx.Where("department_id = ?", deptID).Delete(&department.DepartmentApprovalChain{}).Error; err != nil {
			return err
		}
		// 插入新的
		for i := range chain {
			chain[i].DepartmentID = deptID
		}
		return tx.Create(&chain).Error
	})
}
