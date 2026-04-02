package repository

import (
	"workflow-system/internal/domain/company"

	"gorm.io/gorm"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) Create(company *company.Company) error {
	return r.db.Create(company).Error
}

func (r *CompanyRepository) GetByID(id int64) (*company.Company, error) {
	var comp company.Company
	err := r.db.First(&comp, id).Error
	if err != nil {
		return nil, err
	}
	return &comp, nil
}

func (r *CompanyRepository) GetByCode(code string) (*company.Company, error) {
	var comp company.Company
	err := r.db.Where("code = ?", code).First(&comp).Error
	if err != nil {
		return nil, err
	}
	return &comp, nil
}

func (r *CompanyRepository) List() ([]company.Company, error) {
	var companies []company.Company
	err := r.db.Find(&companies).Error
	return companies, err
}

func (r *CompanyRepository) Update(comp *company.Company) error {
	return r.db.Save(comp).Error
}

// UpdateFields 更新指定字段（支持显式 null 值）
func (r *CompanyRepository) UpdateFields(id int64, fields map[string]interface{}) error {
	return r.db.Model(&company.Company{}).Where("id = ?", id).Updates(fields).Error
}

func (r *CompanyRepository) Delete(id int64) error {
	return r.db.Model(&company.Company{}).Where("id = ?", id).Update("status", 2).Error
}
