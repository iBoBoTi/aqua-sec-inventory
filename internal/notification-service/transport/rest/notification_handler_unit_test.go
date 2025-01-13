package rest_test

import (
	//"bytes"
	"encoding/json"
	//"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/iBoBoTi/aqua-sec-inventory/internal/notification-service/domain"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/notification-service/transport/rest"
)

// Mock NotifiicationUsecase
type mockNotificationUsecase struct {
	mock.Mock
}

func (m *mockNotificationUsecase) CreateNotification(userID int64, message string) (*domain.Notification, error) {
	args := m.Called(userID, message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Notification), args.Error(1)
}
func (m *mockNotificationUsecase) GetAllNotifications(userID int64) ([]domain.Notification, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Notification), args.Error(1)
}
func (m *mockNotificationUsecase) ClearNotification(notificationID int64) error {
	args := m.Called(notificationID)
	return args.Error(0)
}
func (m *mockNotificationUsecase) ClearAllNotifications(userID int64) error {
	args := m.Called(userID)
	return args.Error(0)
}

func TestGetAllUsersNotificationsHandler_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(mockNotificationUsecase)
	handler := rest.NewNotificationHandler(mockUC)

	// Setup Gin
	r := gin.Default()
	r.GET("/users/:id/notifications", handler.GetAllUsersNotifications)

	mockUC.On("GetAllNotifications", int64(1)).Return([]domain.Notification{
		{ID: 1,
			UserID:  1,
			Message: "ebuka",
		},
	}, nil)

	req, _ := http.NewRequest("GET", "/users/1/notifications", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	data, ok := resp["data"].([]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, data)

	mockUC.AssertExpectations(t)
}

func TestGetAllUsersNotificationsHandler_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(mockNotificationUsecase)
	handler := rest.NewNotificationHandler(mockUC)

	// Setup Gin
	r := gin.Default()
	r.GET("/users/:id/notifications", handler.GetAllUsersNotifications)

	req, _ := http.NewRequest("GET", "/users/abc/notifications", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "invalid user id", resp["error"])

	mockUC.AssertExpectations(t)
}

func TestClearSingleNotificationHandler_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(mockNotificationUsecase)
	handler := rest.NewNotificationHandler(mockUC)

	// Setup Gin
	r := gin.Default()
	r.DELETE("/notifications/:id", handler.ClearSingleNotification)

	mockUC.On("ClearNotification", int64(2)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/notifications/2", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Notification cleared", resp["message"])

	mockUC.AssertExpectations(t)
}

func TestClearSingleNotificationHandler_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(mockNotificationUsecase)
	handler := rest.NewNotificationHandler(mockUC)

	// Setup Gin
	r := gin.Default()
	r.DELETE("/notifications/:id", handler.ClearSingleNotification)

	req, _ := http.NewRequest("DELETE", "/notifications/abc", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "invalid notification id", resp["error"])

	mockUC.AssertExpectations(t)
}

func TestClearAllUserNotificationsHandler_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(mockNotificationUsecase)
	handler := rest.NewNotificationHandler(mockUC)

	// Setup Gin
	r := gin.Default()
	r.DELETE("/users/:id/notifications", handler.ClearAllUsersNotifications)

	mockUC.On("ClearAllNotifications", int64(2)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/users/2/notifications", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "All notifications cleared", resp["message"])

	mockUC.AssertExpectations(t)
}

func TestClearAllUserNotificationsHandler_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(mockNotificationUsecase)
	handler := rest.NewNotificationHandler(mockUC)

	// Setup Gin
	r := gin.Default()
	r.DELETE("/users/:id/notifications", handler.ClearAllUsersNotifications)

	req, _ := http.NewRequest("DELETE", "/users/abc/notifications", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "invalid user id", resp["error"])

	mockUC.AssertExpectations(t)
}
