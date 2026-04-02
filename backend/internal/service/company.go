package service

import (
	"workflow-system/internal/domain/company"
	"workflow-system/internal/repository"
)

type CompanyService struct {
	repo *repository.CompanyRepository
}

func NewCompanyService(repo *repository.CompanyRepository) *CompanyService {
	return &CompanyService{repo: repo}
}

func (s *CompanyService) Create(comp *company.Company) error {
	return s.repo.Create(comp)
}

func (s *CompanyService) GetByID(id int64) (*company.Company, error) {
	return s.repo.GetByID(id)
}

func (s *CompanyService) List() ([]company.Company, error) {
	return s.repo.List()
}

func (s *CompanyService) Update(comp *company.Company) error {
	return s.repo.Update(comp)
}

func (s *CompanyService) Delete(id int64) error {
	return s.repo.Delete(id)
}
