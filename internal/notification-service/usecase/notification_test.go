package usecase_test

import (
	// "errors"
	"testing"

	"github.com/iBoBoTi/aqua-sec-inventory/internal/notification-service/domain"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/notification-service/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository
type mockNotificationRepo struct {
	mock.Mock
}

func (m *mockNotificationRepo) Create(notification *domain.Notification) error {
	args := m.Called(notification)
	return args.Error(0)
}
func (m *mockNotificationRepo) GetAllByUserID(userID int64) ([]domain.Notification, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Notification), args.Error(1)
}
func (m *mockNotificationRepo) DeleteByID(notificationID int64) error {
	args := m.Called(notificationID)
	return args.Error(0)
}
func (m *mockNotificationRepo) DeleteAllByUserID(userID int64) error {
	args := m.Called(userID)
	return args.Error(0)
}

func TestCreateNotification_OK(t *testing.T) {
	repo := new(mockNotificationRepo)
	uc := usecase.NewNotificationUsecase(repo)

	repo.On("Create", mock.AnythingOfType("*domain.Notification")).Return(nil)

	cust, err := uc.CreateNotification(int64(2), "test message")
	assert.NoError(t, err)
	assert.NotNil(t, cust)

	repo.AssertExpectations(t)
}

func TestCreateNotification_InvalidID(t *testing.T) {
	repo := new(mockNotificationRepo)
	uc := usecase.NewNotificationUsecase(repo)

	cust, err := uc.CreateNotification(int64(0), "test message")
	assert.EqualError(t, err, "invalid user id")
	assert.Nil(t, cust)

	repo.AssertExpectations(t)
}

func TestCreateNotification_EmptyMessage(t *testing.T) {
	repo := new(mockNotificationRepo)
	uc := usecase.NewNotificationUsecase(repo)

	cust, err := uc.CreateNotification(int64(2), "")
	assert.EqualError(t, err, "empty notification message")
	assert.Nil(t, cust)

	repo.AssertExpectations(t)
}

func TestGetAllNotifications_OK(t *testing.T) {
	repo := new(mockNotificationRepo)
	uc := usecase.NewNotificationUsecase(repo)

	repo.On("GetAllByUserID", int64(1)).Return([]domain.Notification{
		{
			ID: 1,
		},
	}, nil)

	notification, err := uc.GetAllNotifications(int64(1))
	assert.NoError(t, err)
	assert.NotNil(t, notification)

	repo.AssertExpectations(t)
}

func TestGetAllNotifications_InvalidID(t *testing.T) {
	repo := new(mockNotificationRepo)
	uc := usecase.NewNotificationUsecase(repo)

	cust, err := uc.GetAllNotifications(int64(0))
	assert.EqualError(t, err, "invalid user id")
	assert.Nil(t, cust)

	repo.AssertExpectations(t)
}

func TestClearNotification_OK(t *testing.T) {
	repo := new(mockNotificationRepo)
	uc := usecase.NewNotificationUsecase(repo)

	repo.On("DeleteByID", int64(1)).Return(nil)

	err := uc.ClearNotification(int64(1))
	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestClearNotification_InvalidID(t *testing.T) {
	repo := new(mockNotificationRepo)
	uc := usecase.NewNotificationUsecase(repo)

	err := uc.ClearNotification(int64(0))
	assert.EqualError(t, err, "invalid notification id")

	repo.AssertExpectations(t)
}

func TestClearAllNotifications_OK(t *testing.T) {
	repo := new(mockNotificationRepo)
	uc := usecase.NewNotificationUsecase(repo)

	repo.On("DeleteAllByUserID", int64(1)).Return(nil)

	err := uc.ClearAllNotifications(int64(1))
	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestClearAllNotifications_InvalidID(t *testing.T) {
	repo := new(mockNotificationRepo)
	uc := usecase.NewNotificationUsecase(repo)

	err := uc.ClearAllNotifications(int64(0))
	assert.EqualError(t, err, "invalid user id")

	repo.AssertExpectations(t)
}
