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

// GET /notifications/:user_id
func (h *NotificationHandler) GetAll(c *gin.Context) {
    userIDParam := c.Param("user_id")
    userID, err := strconv.ParseInt(userIDParam, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
        return
    }

    notifs, err := h.notificationUC.GetAllNotifications(userID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, notifs)
}

// DELETE /notifications/:user_id
func (h *NotificationHandler) ClearAll(c *gin.Context) {
    userIDParam := c.Param("user_id")
    userID, err := strconv.ParseInt(userIDParam, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
        return
    }

    err = h.notificationUC.ClearAllNotifications(userID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "All notifications cleared"})
}

// DELETE /notifications/:user_id/:notification_id
func (h *NotificationHandler) ClearSingle(c *gin.Context) {
    notifIDParam := c.Param("notification_id")
    notifID, err := strconv.ParseInt(notifIDParam, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification_id"})
        return
    }

    err = h.notificationUC.ClearNotification(notifID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Notification cleared"})
}
