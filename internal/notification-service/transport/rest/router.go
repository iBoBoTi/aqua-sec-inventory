package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/notification-service/service"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/notification-service/usecase"
)

func NewRouter(
    notificationUC usecase.NotificationUsecase,
    notifier service.Notifier, // if you want to use it in the handlers
) *gin.Engine {
    r := gin.Default()


	apiRouter := r.Group("/api/v1/")

    // Notification endpoints
    notificationHandler := NewNotificationHandler(notificationUC)
    apiRouter.GET("/users/:id/notifications", notificationHandler.GetAllUsersNotifications)
    apiRouter.DELETE("/users/:id/notifications", notificationHandler.ClearAllUsersNotifications)
    apiRouter.DELETE("/notifications/:id", notificationHandler.ClearSingleNotification)
   

    return r
}
