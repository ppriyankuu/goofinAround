package services

import (
	"api-testing/interfaces"
	"api-testing/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

var _ interfaces.DBInterface = (*MockDB)(nil) // Ensure MockDB implements DBInterface

func TestCreateUser(t *testing.T) {
	mockDB := new(MockDB)
	userService := NewUserService(mockDB)

	user := &models.User{Name: "John Doe", Email: "john@example.com"}
	mockDB.On("Create", user).Return(&gorm.DB{})

	err := userService.CreateUser(user)
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestGetUserByID(t *testing.T) {
	mockDB := new(MockDB)
	userService := NewUserService(mockDB)

	user := models.User{ID: 1, Name: "John Doe", Email: "john@example.com"}
	mockDB.On("First", &user, uint(1)).Return(&gorm.DB{})

	retrievedUser, err := userService.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, user, *retrievedUser)
	mockDB.AssertExpectations(t)

	// Test not found case
	mockDB.On("First", mock.Anything, uint(2)).Return(&gorm.DB{Error: errors.New("record not found")})
	_, err = userService.GetUserByID(2)
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
}
