package rest_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/iBoBoTi/aqua-sec-inventory/internal/notification-service/repository"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/notification-service/transport/rest"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/notification-service/usecase"
)

func TestGetAllUsersNotificationsHandler_IntegrationTest_OK(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := setUpTestDB(t, "testdb", "testuser", "testpassword")
	defer db.Close()
	assert.NoError(t, err)
	createdNotification := seedNotification(t, db)

	r := gin.Default()
	repo := repository.NewNotificationRepository(db)
	resourceUC := usecase.NewNotificationUsecase(repo)
	handler := rest.NewNotificationHandler(resourceUC)

	r.GET("/users/:id/notifications", handler.GetAllUsersNotifications)

	// Perform the test request
	url := fmt.Sprintf("/users/%d/notifications", createdNotification.UserID)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	notifications := response["data"].([]interface{})

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 1, len(notifications))

}

func TestGetAllUsersNotificationsHandler_IntegrationTest_InvalidUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := setUpTestDB(t, "testdb", "testuser", "testpassword")
	defer db.Close()
	assert.NoError(t, err)

	r := gin.Default()
	repo := repository.NewNotificationRepository(db)
	resourceUC := usecase.NewNotificationUsecase(repo)
	handler := rest.NewNotificationHandler(resourceUC)

	r.GET("/users/:id/notifications", handler.GetAllUsersNotifications)

	// Perform the test request
	req, _ := http.NewRequest(http.MethodGet, "/users/abc/notifications", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "invalid user id", response["error"])

}

func TestClearAllUsersNotificationsHandler_IntegrationTest_OK(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := setUpTestDB(t, "testdb", "testuser", "testpassword")
	defer db.Close()
	assert.NoError(t, err)

	createdNotification := seedNotification(t, db)

	r := gin.Default()
	repo := repository.NewNotificationRepository(db)
	resourceUC := usecase.NewNotificationUsecase(repo)
	handler := rest.NewNotificationHandler(resourceUC)

	r.DELETE("/users/:id/notifications", handler.ClearAllUsersNotifications)

	// Perform the test request
	url := fmt.Sprintf("/users/%d/notifications", createdNotification.UserID)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assertions
	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "All notifications cleared", response["message"])

	notifs, err := repo.GetAllByUserID(createdNotification.UserID)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(notifs))

}

func TestClearAllUsersNotificationsHandler_IntegrationTest_InvalidUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := setUpTestDB(t, "testdb", "testuser", "testpassword")
	defer db.Close()
	assert.NoError(t, err)

	r := gin.Default()
	repo := repository.NewNotificationRepository(db)
	resourceUC := usecase.NewNotificationUsecase(repo)
	handler := rest.NewNotificationHandler(resourceUC)

	r.DELETE("/users/:id/notifications", handler.ClearAllUsersNotifications)

	// Perform the test request
	req, _ := http.NewRequest(http.MethodDelete, "/users/abc/notifications", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "invalid user id", response["error"])

}

func TestClearSingleNotificationsHandler_IntegrationTest__OK(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := setUpTestDB(t, "testdb", "testuser", "testpassword")
	defer db.Close()
	assert.NoError(t, err)
	createdNotification := seedNotification(t, db)

	r := gin.Default()
	repo := repository.NewNotificationRepository(db)
	resourceUC := usecase.NewNotificationUsecase(repo)
	handler := rest.NewNotificationHandler(resourceUC)

	r.DELETE("/notifications/:id", handler.ClearSingleNotification)

	// Perform the test request
	url := fmt.Sprintf("/notifications/%d", createdNotification.ID)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Notification cleared", response["message"])
}

func TestClearSingleNotificationsHandler_IntegrationTest_InvalidUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := setUpTestDB(t, "testdb", "testuser", "testpassword")
	defer db.Close()
	assert.NoError(t, err)

	r := gin.Default()
	repo := repository.NewNotificationRepository(db)
	resourceUC := usecase.NewNotificationUsecase(repo)
	handler := rest.NewNotificationHandler(resourceUC)

	r.DELETE("/notifications/:id", handler.ClearSingleNotification)

	// Perform the test request
	req, _ := http.NewRequest(http.MethodDelete, "/notifications/abc", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "invalid notification id", response["error"])

}
