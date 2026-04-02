package repository

import (
	"workflow-system/internal/domain/position"

	"gorm.io/gorm"
)

type PositionRepository struct {
	db *gorm.DB
}

func NewPositionRepository(db *gorm.DB) *PositionRepository {
	return &PositionRepository{db: db}
}

func (r *PositionRepository) Create(p *position.Position) error {
	return r.db.Create(p).Error
}

func (r *PositionRepository) GetByID(id int64) (*position.Position, error) {
	var p position.Position
	err := r.db.First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PositionRepository) List(companyID int64) ([]position.Position, error) {
	var positions []position.Position
	err := r.db.Where("company_id = ?", companyID).Order("id").Find(&positions).Error
	return positions, err
}

func (r *PositionRepository) Update(p *position.Position) error {
	updates := map[string]interface{}{}
	if p.Name != "" {
		updates["name"] = p.Name
	}
	if p.Code != "" {
		updates["code"] = p.Code
	}
	if p.Status != 0 {
		updates["status"] = p.Status
	}
	return r.db.Model(&position.Position{}).Where("id = ?", p.ID).Updates(updates).Error
}

func (r *PositionRepository) Delete(id int64) error {
	return r.db.Delete(&position.Position{}, id).Error
}

func (r *PositionRepository) GetByLevel(companyID int64, level string) (*position.Position, error) {
	var p position.Position
	// level 字段存的是 name，直接匹配
	err := r.db.Where("company_id = ? AND LOWER(name) = LOWER(?)", companyID, level).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}
