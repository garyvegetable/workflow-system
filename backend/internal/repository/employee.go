package repository

import (
	"workflow-system/internal/domain/employee"

	"gorm.io/gorm"
)

type EmployeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) Create(emp *employee.Employee) error {
	return r.db.Create(emp).Error
}

func (r *EmployeeRepository) GetByID(id int64) (*employee.Employee, error) {
	var emp employee.Employee
	err := r.db.First(&emp, id).Error
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *EmployeeRepository) GetByUsername(username string) (*employee.Employee, error) {
	var emp employee.Employee
	err := r.db.Where("username = ?", username).First(&emp).Error
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *EmployeeRepository) List(companyID int64) ([]employee.Employee, error) {
	var employees []employee.Employee
	err := r.db.Where("company_id = ?", companyID).Find(&employees).Error
	return employees, err
}

func (r *EmployeeRepository) Update(emp *employee.Employee) error {
	return r.db.Save(emp).Error
}

func (r *EmployeeRepository) Delete(id int64) error {
	return r.db.Delete(&employee.Employee{}, id).Error
}

func (r *EmployeeRepository) AddToDepartment(empID, deptID int64) error {
	return r.db.Exec("INSERT INTO employee_department (employee_id, department_id) VALUES (?, ?)", empID, deptID).Error
}

func (r *EmployeeRepository) RemoveFromDepartment(empID, deptID int64) error {
	return r.db.Exec("DELETE FROM employee_department WHERE employee_id = ? AND department_id = ?", empID, deptID).Error
}

func (r *EmployeeRepository) GetDepartments(empID int64) ([]int64, error) {
	var deptIDs []int64
	err := r.db.Table("employee_department").Where("employee_id = ?", empID).Pluck("department_id", &deptIDs).Error
	return deptIDs, err
}

// GetByIDWithDepartments retrieves an employee with their department IDs
func (r *EmployeeRepository) GetByIDWithDepartments(id int64) (*employee.Employee, []int64, error) {
	emp, err := r.GetByID(id)
	if err != nil {
		return nil, nil, err
	}
	deptIDs, err := r.GetDepartments(id)
	if err != nil {
		return nil, nil, err
	}
	return emp, deptIDs, nil
}

// SetDepartments sets the departments for an employee (replaces all existing associations)
func (r *EmployeeRepository) SetDepartments(empID int64, deptIDs []int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Remove all existing associations
		if err := tx.Where("employee_id = ?", empID).Delete(&employee.EmployeeDepartment{}).Error; err != nil {
			return err
		}
		// Add new associations
		for _, deptID := range deptIDs {
			if err := tx.Create(&employee.EmployeeDepartment{
				EmployeeID:   empID,
				DepartmentID: deptID,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ListByDepartmentID retrieves employees belonging to a specific department
func (r *EmployeeRepository) ListByDepartmentID(deptID int64) ([]employee.Employee, error) {
	var employees []employee.Employee
	err := r.db.Table("employee_department ed").
		Select("e.*").
		Joins("LEFT JOIN employee e ON e.id = ed.employee_id").
		Where("ed.department_id = ?", deptID).
		Find(&employees).Error
	return employees, err
}

// ListByCompanyID retrieves all employees for a company
func (r *EmployeeRepository) ListByCompanyID(companyID int64) ([]employee.Employee, error) {
	var employees []employee.Employee
	err := r.db.Where("company_id = ?", companyID).Find(&employees).Error
	return employees, err
}

// SearchByName searches employees by name (case-insensitive partial match)
func (r *EmployeeRepository) SearchByName(name string, companyID int64) ([]employee.Employee, error) {
	var employees []employee.Employee
	err := r.db.Where("company_id = ? AND (name LIKE ? OR username LIKE ?)", companyID, "%"+name+"%", "%"+name+"%").
		Limit(20).
		Find(&employees).Error
	return employees, err
}

// GetCompanyIDByEmployeeID returns the company_id for an employee
func (r *EmployeeRepository) GetCompanyIDByEmployeeID(empID int64) (int64, error) {
	var emp employee.Employee
	err := r.db.Select("company_id").First(&emp, empID).Error
	if err != nil {
		return 0, err
	}
	return emp.CompanyID, nil
}
