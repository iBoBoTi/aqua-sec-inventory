package rest

import (
    "github.com/gin-gonic/gin"
    "github.com/iBoBoTi/aqua-sec-inventory/internal/service"
    "github.com/iBoBoTi/aqua-sec-inventory/internal/usecase"
)

func NewRouter(
    customerUC usecase.CustomerUsecase,
    resourceUC usecase.ResourceUsecase,
    notificationUC usecase.NotificationUsecase,
    notifier service.Notifier, // if you want to use it in the handlers
) *gin.Engine {
    r := gin.Default()


	apiRouter := r.Group("/api/v1/")

    // Customer endpoints
    customerHandler := NewCustomerHandler(customerUC)
    apiRouter.POST("/customers", customerHandler.CreateCustomer)
    apiRouter.GET("/customers/:id", customerHandler.GetCustomerByID)

    // Resource endpoints
    resourceHandler := NewResourceHandler(resourceUC)
    apiRouter.POST("/customers/:id/resources", resourceHandler.AddCloudResource)
    apiRouter.GET("/customers/:id/resources", resourceHandler.GetResourcesByCustomer)
	apiRouter.GET("/resources", resourceHandler.GetAllAvailableResources)
    apiRouter.PUT("/resources/:id", resourceHandler.UpdateResource)
    apiRouter.DELETE("/resources/:id", resourceHandler.DeleteResource)

    // Notification endpoints
    notificationHandler := NewNotificationHandler(notificationUC)
    apiRouter.GET("/notifications/:user_id", notificationHandler.GetAll)
    apiRouter.DELETE("/notifications/:user_id", notificationHandler.ClearAll)
    apiRouter.DELETE("/notifications/:user_id/:notification_id", notificationHandler.ClearSingle)

    return r
}
