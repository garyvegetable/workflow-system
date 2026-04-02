package expense_category

type ExpenseCategory struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyID int64  `gorm:"not null" json:"company_id"`
	Code      string `gorm:"type:varchar(50);not null" json:"code"`
	Name      string `gorm:"type:varchar(200);not null" json:"name"`
	ParentID  *int64 `gorm:"default:null" json:"parent_id"`
	Status    int16  `gorm:"default:1" json:"status"`
}

func (ExpenseCategory) TableName() string {
	return "expense_category"
}

// TreeNode 树形费用科目节点
type ExpenseCategoryTreeNode struct {
	ID        int64                       `json:"id"`
	CompanyID int64                       `json:"company_id"`
	Code      string                      `json:"code"`
	Name      string                      `json:"name"`
	ParentID  *int64                      `json:"parent_id"`
	Status    int16                       `json:"status"`
	Children  []*ExpenseCategoryTreeNode  `json:"children,omitempty"`
}
