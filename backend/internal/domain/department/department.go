package department

type Department struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyID int64  `gorm:"not null" json:"company_id"`
	ParentID  *int64 `gorm:"default:null" json:"parent_id"`
	Name      string `gorm:"type:varchar(100);not null" json:"name"`
	LeaderID  *int64 `json:"leader_id"`
	Status    int16  `gorm:"default:1" json:"status"`
}

func (Department) TableName() string {
	return "department"
}

// TreeNode 树形部门节点
type TreeNode struct {
	ID        int64       `json:"id"`
	CompanyID int64       `json:"company_id"`
	ParentID  *int64      `json:"parent_id"`
	Name      string      `json:"name"`
	LeaderID  *int64      `json:"leader_id"`
	Status    int16       `json:"status"`
	Children  []*TreeNode `json:"children,omitempty"`
}
