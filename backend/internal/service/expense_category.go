package service

import (
	"workflow-system/internal/domain/expense_category"
	"workflow-system/internal/repository"
)

type ExpenseCategoryService struct {
	repo *repository.ExpenseCategoryRepository
}

func NewExpenseCategoryService(repo *repository.ExpenseCategoryRepository) *ExpenseCategoryService {
	return &ExpenseCategoryService{repo: repo}
}

func (s *ExpenseCategoryService) Create(cat *expense_category.ExpenseCategory) error {
	return s.repo.Create(cat)
}

func (s *ExpenseCategoryService) GetByID(id int64) (*expense_category.ExpenseCategory, error) {
	return s.repo.GetByID(id)
}

func (s *ExpenseCategoryService) List(companyID int64) ([]expense_category.ExpenseCategory, error) {
	return s.repo.List(companyID)
}

func (s *ExpenseCategoryService) GetTree(companyID int64) ([]*expense_category.ExpenseCategoryTreeNode, error) {
	return s.repo.GetTree(companyID)
}

func (s *ExpenseCategoryService) Update(cat *expense_category.ExpenseCategory) error {
	return s.repo.Update(cat)
}

func (s *ExpenseCategoryService) Delete(id int64) error {
	return s.repo.Delete(id)
}
