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

	"workflow-system/internal/domain/position"
	"workflow-system/internal/service"
)

// Mock PositionService
type MockPositionService struct {
	mock.Mock
}

func (m *MockPositionService) Create(pos *position.Position) error {
	args := m.Called(pos)
	return args.Error(0)
}

func (m *MockPositionService) GetByID(id int64) (*position.Position, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*position.Position), args.Error(1)
}

func (m *MockPositionService) List(companyID int64) ([]position.Position, error) {
	args := m.Called(companyID)
	return args.Get(0).([]position.Position), args.Error(1)
}

func (m *MockPositionService) Update(pos *position.Position) error {
	args := m.Called(pos)
	return args.Error(0)
}

func (m *MockPositionService) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupPositionRouter(handler *PositionHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/positions", handler.List)
	r.POST("/positions", handler.Create)
	r.GET("/positions/:id", handler.Get)
	r.PUT("/positions/:id", handler.Update)
	r.DELETE("/positions/:id", handler.Delete)
	return r
}

func TestPositionHandler_List(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)
	router := setupPositionRouter(handler)

	expectedPositions := []position.Position{
		{ID: 1, CompanyID: 1, Name: "Engineer", Status: 1},
		{ID: 2, CompanyID: 1, Name: "Manager", Status: 1},
	}

	mockService.On("List", int64(1)).Return(expectedPositions, nil)

	req, _ := http.NewRequest("GET", "/positions?company_id=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var positions []position.Position
	err := json.Unmarshal(w.Body.Bytes(), &positions)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(positions))
	mockService.AssertExpectations(t)
}

func TestPositionHandler_Create(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)
	router := setupPositionRouter(handler)

	newPos := position.Position{
		Name: "Director",
	}

	mockService.On("Create", mock.AnythingOfType("*position.Position")).Return(nil)

	body, _ := json.Marshal(newPos)
	req, _ := http.NewRequest("POST", "/positions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestPositionHandler_Get(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)
	router := setupPositionRouter(handler)

	expectedPos := &position.Position{
		ID:       1,
		CompanyID: 1,
		Name:     "Engineer",
		Status:   1,
	}

	mockService.On("GetByID", int64(1)).Return(expectedPos, nil)

	req, _ := http.NewRequest("GET", "/positions/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var pos position.Position
	err := json.Unmarshal(w.Body.Bytes(), &pos)
	assert.NoError(t, err)
	assert.Equal(t, "Engineer", pos.Name)
	mockService.AssertExpectations(t)
}

func TestPositionHandler_Get_NotFound(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)
	router := setupPositionRouter(handler)

	mockService.On("GetByID", int64(999)).Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/positions/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestPositionHandler_Update(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)
	router := setupPositionRouter(handler)

	updatedPos := position.Position{
		Name: "Senior Engineer",
	}

	mockService.On("Update", mock.AnythingOfType("*position.Position")).Return(nil)

	body, _ := json.Marshal(updatedPos)
	req, _ := http.NewRequest("PUT", "/positions/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestPositionHandler_Delete(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)
	router := setupPositionRouter(handler)

	mockService.On("Delete", int64(1)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/positions/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
