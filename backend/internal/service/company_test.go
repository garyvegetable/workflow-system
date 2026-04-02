package service

import (
	"testing"
	"workflow-system/internal/domain/company"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock CompanyRepository
type MockCompanyRepository struct {
	mock.Mock
}

func (m *MockCompanyRepository) Create(c *company.Company) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockCompanyRepository) GetByID(id int64) (*company.Company, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*company.Company), args.Error(1)
}

func (m *MockCompanyRepository) GetByCode(code string) (*company.Company, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*company.Company), args.Error(1)
}

func (m *MockCompanyRepository) List() ([]company.Company, error) {
	args := m.Called()
	return args.Get(0).([]company.Company), args.Error(1)
}

func (m *MockCompanyRepository) Update(c *company.Company) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockCompanyRepository) UpdateFields(id int64, fields map[string]interface{}) error {
	args := m.Called(id, fields)
	return args.Error(0)
}

func (m *MockCompanyRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCompanyService_GetByID(t *testing.T) {
	mockRepo := new(MockCompanyRepository)
	svc := NewCompanyService(mockRepo)
	assert.NotNil(t, svc)
}

func TestCompanyService_GetByCode(t *testing.T) {
	mockRepo := new(MockCompanyRepository)
	svc := NewCompanyService(mockRepo)

	expectedCompany := &company.Company{
		ID:       1,
		Code:     "DEMO",
		Name:     "Demo Company",
		Status:   1,
	}

	mockRepo.On("GetByCode", "DEMO").Return(expectedCompany, nil)

	result, err := svc.GetByCode("DEMO")

	assert.NoError(t, err)
	assert.Equal(t, "DEMO", result.Code)
	assert.Equal(t, "Demo Company", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestCompanyService_List(t *testing.T) {
	mockRepo := new(MockCompanyRepository)
	svc := NewCompanyService(mockRepo)

	expectedCompanies := []company.Company{
		{ID: 1, Code: "DEMO", Name: "Demo Company", Status: 1},
		{ID: 2, Code: "TEST", Name: "Test Company", Status: 1},
	}

	mockRepo.On("List").Return(expectedCompanies, nil)

	companies, err := svc.List()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(companies))
	mockRepo.AssertExpectations(t)
}

func TestCompanyService_Create(t *testing.T) {
	mockRepo := new(MockCompanyRepository)
	svc := NewCompanyService(mockRepo)

	newCompany := &company.Company{
		Code:   "NEW",
		Name:   "New Company",
		Status: 1,
	}

	mockRepo.On("Create", newCompany).Return(nil)

	err := svc.Create(newCompany)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCompanyService_UpdateFields(t *testing.T) {
	mockRepo := new(MockCompanyRepository)
	svc := NewCompanyService(mockRepo)

	fields := map[string]interface{}{
		"name": "Updated Company Name",
	}

	mockRepo.On("UpdateFields", int64(1), fields).Return(nil)

	err := svc.UpdateFields(1, fields)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCompanyService_Delete(t *testing.T) {
	mockRepo := new(MockCompanyRepository)
	svc := NewCompanyService(mockRepo)

	mockRepo.On("Delete", int64(1)).Return(nil)

	err := svc.Delete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
