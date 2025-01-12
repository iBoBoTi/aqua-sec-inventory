package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/main-service/usecase"
)

type CustomerHandler struct {
    customerUC usecase.CustomerUsecase
}

func NewCustomerHandler(customerUC usecase.CustomerUsecase) *CustomerHandler {
    return &CustomerHandler{customerUC: customerUC}
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
    var req struct {
        Name  string `json:"name" binding:"required"`
        Email string `json:"email" binding:"required,email"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    customer, err := h.customerUC.CreateCustomer(req.Name, req.Email)
    if err != nil {
		if err.Error() == "internal server error" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
        c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
		"data": map[string]interface{}{
			"id":    customer.ID,
			"name":  customer.Name,
			"email": customer.Email,
		},
    })
}

func (h *CustomerHandler) GetCustomerByID(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.ParseInt(idParam, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
        return
    }

    customer, err := h.customerUC.GetCustomerByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
		"data": map[string]interface{}{
        "id":    customer.ID,
        "name":  customer.Name,
        "email": customer.Email,
    }})
}
