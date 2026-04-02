package service

import (
	"workflow-system/internal/domain/position"
	"workflow-system/internal/repository"
)

type PositionService struct {
	repo *repository.PositionRepository
}

func NewPositionService(repo *repository.PositionRepository) *PositionService {
	return &PositionService{repo: repo}
}

func (s *PositionService) Create(p *position.Position) error {
	return s.repo.Create(p)
}

func (s *PositionService) GetByID(id int64) (*position.Position, error) {
	return s.repo.GetByID(id)
}

func (s *PositionService) List(companyID int64) ([]position.Position, error) {
	return s.repo.List(companyID)
}

func (s *PositionService) Update(p *position.Position) error {
	return s.repo.Update(p)
}

func (s *PositionService) Delete(id int64) error {
	return s.repo.Delete(id)
}
