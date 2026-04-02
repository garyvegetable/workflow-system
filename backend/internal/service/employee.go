package service

import (
	"workflow-system/internal/domain/employee"
	"workflow-system/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type EmployeeService struct {
	repo     *repository.EmployeeRepository
	deptRepo *repository.DepartmentRepository
	bankRepo *repository.EmployeeBankAccountRepository
}

func NewEmployeeService(repo *repository.EmployeeRepository, deptRepo *repository.DepartmentRepository, bankRepo *repository.EmployeeBankAccountRepository) *EmployeeService {
	return &EmployeeService{repo: repo, deptRepo: deptRepo, bankRepo: bankRepo}
}

func (s *EmployeeService) Create(emp *employee.Employee) error {
	// Hash password if provided and not already hashed
	if emp.PasswordHash != "" && len(emp.PasswordHash) < 60 {
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(emp.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		emp.PasswordHash = string(hashedBytes)
	}
	return s.repo.Create(emp)
}

func (s *EmployeeService) GetByID(id int64) (*employee.Employee, error) {
	return s.repo.GetByID(id)
}

func (s *EmployeeService) GetByUsername(username string) (*employee.Employee, error) {
	return s.repo.GetByUsername(username)
}

func (s *EmployeeService) VerifyPassword(emp *employee.Employee, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(emp.PasswordHash), []byte(password))
	return err == nil
}

func (s *EmployeeService) List(companyID int64) ([]employee.Employee, error) {
	return s.repo.List(companyID)
}

func (s *EmployeeService) Update(emp *employee.Employee) error {
	// Hash password if provided and not already hashed (plain text password < 60 chars)
	if emp.PasswordHash != "" && len(emp.PasswordHash) < 60 {
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(emp.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		emp.PasswordHash = string(hashedBytes)
	}
	return s.repo.Update(emp)
}

func (s *EmployeeService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *EmployeeService) ListBankAccounts(empID int64) ([]employee.EmployeeBankAccount, error) {
	return s.bankRepo.ListByEmployeeID(empID)
}

func (s *EmployeeService) CreateBankAccount(account *employee.EmployeeBankAccount) error {
	return s.bankRepo.Create(account)
}

func (s *EmployeeService) UpdateBankAccount(account *employee.EmployeeBankAccount) error {
	return s.bankRepo.Update(account)
}

func (s *EmployeeService) DeleteBankAccount(id int64) error {
	return s.bankRepo.Delete(id)
}

func (s *EmployeeService) SearchByName(name string, companyID int64) ([]employee.Employee, error) {
	return s.repo.SearchByName(name, companyID)
}

func (s *EmployeeService) SetDepartments(empID int64, deptIDs []int64) error {
	return s.repo.SetDepartments(empID, deptIDs)
}

func (s *EmployeeService) GetDepartments(empID int64) ([]int64, error) {
	return s.repo.GetDepartments(empID)
}
