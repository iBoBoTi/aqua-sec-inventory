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
    apiRouter.GET("/notifications/:user_id", notificationHandler.GetAll)
    apiRouter.DELETE("/notifications/:user_id", notificationHandler.ClearAll)
    apiRouter.DELETE("/notifications/:user_id/:notification_id", notificationHandler.ClearSingle)

    return r
}
