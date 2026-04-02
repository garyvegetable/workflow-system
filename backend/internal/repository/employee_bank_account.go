package repository

import (
	"workflow-system/internal/domain/employee"

	"gorm.io/gorm"
)

type EmployeeBankAccountRepository struct {
	db *gorm.DB
}

func NewEmployeeBankAccountRepository(db *gorm.DB) *EmployeeBankAccountRepository {
	return &EmployeeBankAccountRepository{db: db}
}

func (r *EmployeeBankAccountRepository) Create(account *employee.EmployeeBankAccount) error {
	return r.db.Create(account).Error
}

func (r *EmployeeBankAccountRepository) GetByID(id int64) (*employee.EmployeeBankAccount, error) {
	var account employee.EmployeeBankAccount
	err := r.db.First(&account, id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *EmployeeBankAccountRepository) ListByEmployeeID(empID int64) ([]employee.EmployeeBankAccount, error) {
	var accounts []employee.EmployeeBankAccount
	err := r.db.Where("employee_id = ?", empID).Find(&accounts).Error
	return accounts, err
}

func (r *EmployeeBankAccountRepository) Update(account *employee.EmployeeBankAccount) error {
	return r.db.Save(account).Error
}

func (r *EmployeeBankAccountRepository) Delete(id int64) error {
	return r.db.Delete(&employee.EmployeeBankAccount{}, id).Error
}

func (r *EmployeeBankAccountRepository) SetDefault(id int64, empID int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 取消原来的默认
		if err := tx.Model(&employee.EmployeeBankAccount{}).Where("employee_id = ? AND is_default = ?", empID, true).Update("is_default", false).Error; err != nil {
			return err
		}
		// 设置新的默认
		return tx.Model(&employee.EmployeeBankAccount{}).Where("id = ?", id).Update("is_default", true).Error
	})
}
