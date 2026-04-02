package repository

import (
	"workflow-system/internal/domain/expense_category"

	"gorm.io/gorm"
)

type ExpenseCategoryRepository struct {
	db *gorm.DB
}

func NewExpenseCategoryRepository(db *gorm.DB) *ExpenseCategoryRepository {
	return &ExpenseCategoryRepository{db: db}
}

func (r *ExpenseCategoryRepository) Create(category *expense_category.ExpenseCategory) error {
	return r.db.Create(category).Error
}

func (r *ExpenseCategoryRepository) GetByID(id int64) (*expense_category.ExpenseCategory, error) {
	var category expense_category.ExpenseCategory
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *ExpenseCategoryRepository) List(companyID int64) ([]expense_category.ExpenseCategory, error) {
	var categories []expense_category.ExpenseCategory
	err := r.db.Where("company_id = ?", companyID).Find(&categories).Error
	return categories, err
}

func (r *ExpenseCategoryRepository) GetTree(companyID int64) ([]*expense_category.ExpenseCategoryTreeNode, error) {
	var categories []expense_category.ExpenseCategory
	err := r.db.Where("company_id = ?", companyID).Order("id").Find(&categories).Error
	if err != nil {
		return nil, err
	}

	// 构建树形结构
	nodeMap := make(map[int64]*expense_category.ExpenseCategoryTreeNode)
	var roots []*expense_category.ExpenseCategoryTreeNode

	// 先将所有费用科目转为节点
	for i := range categories {
		nodeMap[categories[i].ID] = &expense_category.ExpenseCategoryTreeNode{
			ID:        categories[i].ID,
			CompanyID: categories[i].CompanyID,
			Code:      categories[i].Code,
			Name:      categories[i].Name,
			ParentID:  categories[i].ParentID,
			Status:    categories[i].Status,
			Children:  []*expense_category.ExpenseCategoryTreeNode{},
		}
	}

	// 再构建父子关系
	for i := range categories {
		node := nodeMap[categories[i].ID]
		if categories[i].ParentID != nil && *categories[i].ParentID > 0 {
			if parent, ok := nodeMap[*categories[i].ParentID]; ok {
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

func (r *ExpenseCategoryRepository) Update(category *expense_category.ExpenseCategory) error {
	return r.db.Save(category).Error
}

func (r *ExpenseCategoryRepository) Delete(id int64) error {
	return r.db.Delete(&expense_category.ExpenseCategory{}, id).Error
}
