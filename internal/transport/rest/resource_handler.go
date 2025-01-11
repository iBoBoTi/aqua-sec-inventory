package rest

import (
    "net/http"
    "strconv"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/iBoBoTi/aqua-sec-inventory/internal/usecase"
)

type ResourceHandler struct {
    resourceUC usecase.ResourceUsecase
}

func NewResourceHandler(resourceUC usecase.ResourceUsecase) *ResourceHandler {
    return &ResourceHandler{resourceUC: resourceUC}
}

// GET /resources
func (h *ResourceHandler) GetAllAvailableResources(c *gin.Context) {
	resources, err := h.resourceUC.GetAllAvailableResources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resources})
}

// POST /customers/:id/resources
func (h *ResourceHandler) AddCloudResources(c *gin.Context) {
    customerIDParam := c.Param("id")
    customerID, err := strconv.ParseInt(customerIDParam, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer_id"})
        return
    }

    var req struct {
        ResourceNames []string `json:"resource_names" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Trim whitespace
    for i, name := range req.ResourceNames {
        req.ResourceNames[i] = strings.TrimSpace(name)
    }

    err = h.resourceUC.AddCloudResources(customerID, req.ResourceNames)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Resources assigned successfully"})
}

// POST /customers/:id/resources
func (h *ResourceHandler) AddCloudResource(c *gin.Context) {
    customerIDParam := c.Param("id")
    customerID, err := strconv.ParseInt(customerIDParam, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer id"})
        return
    }

    var req struct {
        ResourceName string `json:"resource_name" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err = h.resourceUC.AddCloudResource(customerID, strings.TrimSpace(req.ResourceName))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Resources assigned successfully"})
}

// GET /customers/:id/resources
func (h *ResourceHandler) GetResourcesByCustomer(c *gin.Context) {
    customerIDParam := c.Param("id")
    customerID, err := strconv.ParseInt(customerIDParam, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer_id"})
        return
    }

    resources, err := h.resourceUC.GetResourcesByCustomer(customerID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, map[string]interface{}{
		"data": resources,
	})
}

// PUT /resources/:id
func (h *ResourceHandler) UpdateResource(c *gin.Context) {
    resourceIDParam := c.Param("id")
    resourceID, err := strconv.ParseInt(resourceIDParam, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource id"})
        return
    }

    var req struct {
        Name       string `json:"name" binding:"required"`
        Type       string `json:"type" binding:"required"`
        Region     string `json:"region" binding:"required"`
        CustomerID *int64 `json:"customer_id"` // optional
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    updatedRes, err := h.resourceUC.UpdateResource(resourceID, req.Name, req.Type, req.Region, req.CustomerID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, updatedRes)
}

// DELETE /resources/:id
func (h *ResourceHandler) DeleteResource(c *gin.Context) {
    resourceIDParam := c.Param("id")
    resourceID, err := strconv.ParseInt(resourceIDParam, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource id"})
        return
    }

    if err := h.resourceUC.DeleteResource(resourceID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Resource deleted successfully"})
}
