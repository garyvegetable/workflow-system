package department

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDepartment_Fields(t *testing.T) {
	parentID := int64(0)
	dept := Department{
		ID:        1,
		CompanyID: 1,
		ParentID:  &parentID,
		Name:      "Engineering",
		LeaderID:  nil,
		Status:    1,
	}

	assert.Equal(t, int64(1), dept.ID)
	assert.Equal(t, int64(1), dept.CompanyID)
	assert.Nil(t, dept.ParentID)
	assert.Equal(t, "Engineering", dept.Name)
	assert.Equal(t, int8(1), dept.Status)
}

func TestTreeNode_Children(t *testing.T) {
	root := &TreeNode{
		ID:       1,
		CompanyID: 1,
		Name:     "Root",
		Children: []*TreeNode{},
	}

	child := &TreeNode{
		ID:       2,
		CompanyID: 1,
		ParentID: ptrInt64(1),
		Name:     "Child",
		Children: []*TreeNode{},
	}

	root.Children = append(root.Children, child)

	assert.Equal(t, 1, len(root.Children))
	assert.Equal(t, "Child", root.Children[0].Name)
}

func TestDepartmentApprovalChain_Fields(t *testing.T) {
	chain := DepartmentApprovalChain{
		ID:           1,
		DepartmentID: 1,
		EmployeeID:   100,
		StepOrder:    1,
	}

	assert.Equal(t, int64(1), chain.ID)
	assert.Equal(t, int64(1), chain.DepartmentID)
	assert.Equal(t, int64(100), chain.EmployeeID)
	assert.Equal(t, 1, chain.StepOrder)
}

func ptrInt64(v int64) *int64 {
	return &v
}
