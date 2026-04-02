package repository

import (
	"workflow-system/internal/domain/supplier"

	"gorm.io/gorm"
)

type SupplierRepository struct {
	db *gorm.DB
}

func NewSupplierRepository(db *gorm.DB) *SupplierRepository {
	return &SupplierRepository{db: db}
}

func (r *SupplierRepository) Create(supplier *supplier.Supplier) error {
	return r.db.Create(supplier).Error
}

func (r *SupplierRepository) GetByID(id int64) (*supplier.Supplier, error) {
	var sup supplier.Supplier
	err := r.db.First(&sup, id).Error
	if err != nil {
		return nil, err
	}
	return &sup, nil
}

func (r *SupplierRepository) List(companyID int64) ([]supplier.Supplier, error) {
	var suppliers []supplier.Supplier
	err := r.db.Where("company_id = ?", companyID).Find(&suppliers).Error
	return suppliers, err
}

func (r *SupplierRepository) Update(supplier *supplier.Supplier) error {
	return r.db.Save(supplier).Error
}

func (r *SupplierRepository) Delete(id int64) error {
	return r.db.Delete(&supplier.Supplier{}, id).Error
}
