package service

import (
	"testing"
	"workflow-system/internal/domain/position"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock PositionRepository
type MockPositionRepository struct {
	mock.Mock
}

func (m *MockPositionRepository) Create(pos *position.Position) error {
	args := m.Called(pos)
	return args.Error(0)
}

func (m *MockPositionRepository) GetByID(id int64) (*position.Position, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*position.Position), args.Error(1)
}

func (m *MockPositionRepository) List(companyID int64) ([]position.Position, error) {
	args := m.Called(companyID)
	return args.Get(0).([]position.Position), args.Error(1)
}

func (m *MockPositionRepository) Update(pos *position.Position) error {
	args := m.Called(pos)
	return args.Error(0)
}

func (m *MockPositionRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestPositionService_GetByID(t *testing.T) {
	mockRepo := new(MockPositionRepository)
	svc := NewPositionService(mockRepo)
	assert.NotNil(t, svc)
}

func TestPositionService_List(t *testing.T) {
	mockRepo := new(MockPositionRepository)
	svc := NewPositionService(mockRepo)

	expectedPositions := []position.Position{
		{ID: 1, CompanyID: 1, Name: "Engineer", Status: 1},
		{ID: 2, CompanyID: 1, Name: "Manager", Status: 1},
	}

	mockRepo.On("List", int64(1)).Return(expectedPositions, nil)

	positions, err := svc.List(1)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(positions))
	assert.Equal(t, "Engineer", positions[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestPositionService_Create(t *testing.T) {
	mockRepo := new(MockPositionRepository)
	svc := NewPositionService(mockRepo)

	newPos := &position.Position{
		CompanyID: 1,
		Name:      "Director",
		Status:    1,
	}

	mockRepo.On("Create", newPos).Return(nil)

	err := svc.Create(newPos)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), newPos.ID) // ID should be set by DB
	mockRepo.AssertExpectations(t)
}

func TestPositionService_Update(t *testing.T) {
	mockRepo := new(MockPositionRepository)
	svc := NewPositionService(mockRepo)

	pos := &position.Position{
		ID:       1,
		CompanyID: 1,
		Name:     "Senior Engineer",
		Status:   1,
	}

	mockRepo.On("Update", pos).Return(nil)

	err := svc.Update(pos)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPositionService_Delete(t *testing.T) {
	mockRepo := new(MockPositionRepository)
	svc := NewPositionService(mockRepo)

	mockRepo.On("Delete", int64(1)).Return(nil)

	err := svc.Delete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
