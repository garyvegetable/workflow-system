package service

import (
	"workflow-system/internal/domain/supplier"
	"workflow-system/internal/repository"
)

type SupplierService struct {
	repo *repository.SupplierRepository
}

func NewSupplierService(repo *repository.SupplierRepository) *SupplierService {
	return &SupplierService{repo: repo}
}

func (s *SupplierService) Create(sup *supplier.Supplier) error {
	return s.repo.Create(sup)
}

func (s *SupplierService) GetByID(id int64) (*supplier.Supplier, error) {
	return s.repo.GetByID(id)
}

func (s *SupplierService) List(companyID int64) ([]supplier.Supplier, error) {
	return s.repo.List(companyID)
}

func (s *SupplierService) Update(sup *supplier.Supplier) error {
	return s.repo.Update(sup)
}

func (s *SupplierService) Delete(id int64) error {
	return s.repo.Delete(id)
}
