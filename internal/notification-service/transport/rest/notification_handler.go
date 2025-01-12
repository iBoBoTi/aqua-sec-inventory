package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/notification-service/usecase"
)

type NotificationHandler struct {
    notificationUC usecase.NotificationUsecase
}

func NewNotificationHandler(notificationUC usecase.NotificationUsecase) *NotificationHandler {
    return &NotificationHandler{notificationUC: notificationUC}
}

// GET /users/:id/notifications
func (h *NotificationHandler) GetAllUsersNotifications(c *gin.Context) {
    userIDParam := c.Param("id")
    userID, err := strconv.ParseInt(userIDParam, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }

    notifications, err := h.notificationUC.GetAllNotifications(userID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, map[string]interface{}{"data": notifications})
}

// DELETE /users/:id/notifications
func (h *NotificationHandler) ClearAllUsersNotifications(c *gin.Context) {
    userIDParam := c.Param("id")
    userID, err := strconv.ParseInt(userIDParam, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }

    err = h.notificationUC.ClearAllNotifications(userID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "All notifications cleared"})
}

// DELETE /notifications/:id
func (h *NotificationHandler) ClearSingleNotification(c *gin.Context) {
    notificationIDParam := c.Param("id")
    notificationID, err := strconv.ParseInt(notificationIDParam, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id"})
        return
    }

    err = h.notificationUC.ClearNotification(notificationID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Notification cleared"})
}
