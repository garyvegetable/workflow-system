package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"workflow-system/internal/domain/company"
	"workflow-system/internal/service"
)

// Mock CompanyService
type MockCompanyService struct {
	mock.Mock
}

func (m *MockCompanyService) Create(c *company.Company) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockCompanyService) GetByID(id int64) (*company.Company, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*company.Company), args.Error(1)
}

func (m *MockCompanyService) GetByCode(code string) (*company.Company, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*company.Company), args.Error(1)
}

func (m *MockCompanyService) List() ([]company.Company, error) {
	args := m.Called()
	return args.Get(0).([]company.Company), args.Error(1)
}

func (m *MockCompanyService) Update(c *company.Company) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockCompanyService) UpdateFields(id int64, fields map[string]interface{}) error {
	args := m.Called(id, fields)
	return args.Error(0)
}

func (m *MockCompanyService) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupCompanyRouter(handler *CompanyHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/companies", handler.List)
	r.POST("/companies", handler.Create)
	r.GET("/companies/:id", handler.Get)
	r.PUT("/companies/:id", handler.Update)
	r.DELETE("/companies/:id", handler.Delete)
	return r
}

func TestCompanyHandler_List(t *testing.T) {
	mockService := new(MockCompanyService)
	handler := NewCompanyHandler(mockService)
	router := setupCompanyRouter(handler)

	expectedCompanies := []company.Company{
		{ID: 1, Code: "DEMO", Name: "Demo Company", Status: 1},
		{ID: 2, Code: "TEST", Name: "Test Company", Status: 1},
	}

	mockService.On("List").Return(expectedCompanies, nil)

	req, _ := http.NewRequest("GET", "/companies", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var companies []company.Company
	err := json.Unmarshal(w.Body.Bytes(), &companies)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(companies))
	mockService.AssertExpectations(t)
}

func TestCompanyHandler_Create(t *testing.T) {
	mockService := new(MockCompanyService)
	handler := NewCompanyHandler(mockService)
	router := setupCompanyRouter(handler)

	newCompany := company.Company{
		Code:   "NEW",
		Name:   "New Company",
		Status: 1,
	}

	mockService.On("Create", mock.AnythingOfType("*company.Company")).Return(nil)

	body, _ := json.Marshal(newCompany)
	req, _ := http.NewRequest("POST", "/companies", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestCompanyHandler_Get(t *testing.T) {
	mockService := new(MockCompanyService)
	handler := NewCompanyHandler(mockService)
	router := setupCompanyRouter(handler)

	expectedCompany := &company.Company{
		ID:       1,
		Code:     "DEMO",
		Name:     "Demo Company",
		Status:   1,
	}

	mockService.On("GetByID", int64(1)).Return(expectedCompany, nil)

	req, _ := http.NewRequest("GET", "/companies/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var c company.Company
	err := json.Unmarshal(w.Body.Bytes(), &c)
	assert.NoError(t, err)
	assert.Equal(t, "DEMO", c.Code)
	mockService.AssertExpectations(t)
}

func TestCompanyHandler_Get_NotFound(t *testing.T) {
	mockService := new(MockCompanyService)
	handler := NewCompanyHandler(mockService)
	router := setupCompanyRouter(handler)

	mockService.On("GetByID", int64(999)).Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/companies/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestCompanyHandler_Update(t *testing.T) {
	mockService := new(MockCompanyService)
	handler := NewCompanyHandler(mockService)
	router := setupCompanyRouter(handler)

	updatedCompany := company.Company{
		Name: "Updated Company",
	}

	mockService.On("UpdateFields", int64(1), mock.AnythingOfType("map[string]interface {}")).Return(nil)
	mockService.On("GetByID", int64(1)).Return(&company.Company{ID: 1, Name: "Updated Company"}, nil)

	body, _ := json.Marshal(updatedCompany)
	req, _ := http.NewRequest("PUT", "/companies/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestCompanyHandler_Delete(t *testing.T) {
	mockService := new(MockCompanyService)
	handler := NewCompanyHandler(mockService)
	router := setupCompanyRouter(handler)

	mockService.On("Delete", int64(1)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/companies/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
