package rest_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/iBoBoTi/aqua-sec-inventory/internal/main-service/domain"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/main-service/transport/rest"
)

// Mock ResourceUsecase
type mockResourceUsecase struct {
	mock.Mock
}

func (m *mockResourceUsecase) GetAllAvailableResources() ([]domain.Resource, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Resource), args.Error(1)
}

func (m *mockResourceUsecase) AddCloudResource(customerID int64, resourceName string) error {
	args := m.Called(customerID, resourceName)
	return args.Error(0)
}

func (m *mockResourceUsecase) AddCloudResources(customerID int64, resourceNames []string) error {
	args := m.Called(customerID, resourceNames)
	return args.Error(0)
}

func (m *mockResourceUsecase) GetResourcesByCustomer(customerID int64) ([]domain.Resource, error) {
	args := m.Called(customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Resource), args.Error(1)
}

func (m *mockResourceUsecase) UpdateResource(resourceID int64, name, resourceType, region string) (*domain.Resource, error) {
	args := m.Called(resourceID, name, resourceType, region)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Resource), args.Error(1)
}

func (m *mockResourceUsecase) DeleteResource(resourceID int64) error {
	args := m.Called(resourceID)
	return args.Error(0)
}

type mockNotifier struct {
	mock.Mock
}

func (m *mockNotifier) Publish(message domain.Notification) error {
	args := m.Called(message)
	return args.Error(0)
}
func (m *mockNotifier) Listen() error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockNotifier) Close() {
	m.Called()
}

func TestAddCloudResourceHandler_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(mockResourceUsecase)
	mockNotify := new(mockNotifier)
	handler := rest.NewResourceHandler(mockUC, mockNotify)

	// Setup Gin
	r := gin.Default()
	r.POST("/customers/:id/resources", handler.AddCloudResource)

	mockUC.On("AddCloudResource", int64(123), "aws_vpc_main").Return(nil)
	mockNotify.On("Publish", domain.Notification{
		Event:   "notification",
		UserID:  int64(123),
		Message: "added resource aws_vpc_main for customer with customerID 123",
	}).Return(nil)

	body := `{"resource_name":"aws_vpc_main"}`
	req, _ := http.NewRequest("POST", "/customers/123/resources", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Resources assigned successfully", resp["message"])

	mockUC.AssertExpectations(t)
}

func TestAddCloudResourceHandler_InvalidCustomerID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mockResourceUsecase)
	mockNotify := new(mockNotifier)
	handler := rest.NewResourceHandler(mockUC, mockNotify)

	r := gin.Default()
	r.POST("/customers/:id/resources", handler.AddCloudResource)

	body := `{"resource_name":"aws_vpc_main"}`
	req, _ := http.NewRequest("POST", "/customers/abc/resources", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "invalid customer id", resp["error"])
}

func TestAddCloudResourcesHandler_CustomerNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(mockResourceUsecase)
	mockNotify := new(mockNotifier)
	handler := rest.NewResourceHandler(mockUC, mockNotify)

	// Setup Gin
	r := gin.Default()
	r.POST("/customers/:id/resources", handler.AddCloudResource)

	mockUC.On("AddCloudResource", int64(123), "aws_vpc_main").Return(errors.New("customer not found"))

	body := `{"resource_name":"aws_vpc_main"}`
	req, _ := http.NewRequest("POST", "/customers/123/resources", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "customer not found", resp["error"])

	mockUC.AssertExpectations(t)
}

func TestAddCloudResourcesHandler_ResourceAlreadyExist(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(mockResourceUsecase)
	mockNotify := new(mockNotifier)
	handler := rest.NewResourceHandler(mockUC, mockNotify)

	// Setup Gin
	r := gin.Default()
	r.POST("/customers/:id/resources", handler.AddCloudResource)

	mockUC.On("AddCloudResource", int64(123), "aws_vpc_main").Return(errors.New("customer already has aws_vpc_main resource"))

	body := `{"resource_name":"aws_vpc_main"}`
	req, _ := http.NewRequest("POST", "/customers/123/resources", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "customer already has aws_vpc_main resource", resp["error"])

	mockUC.AssertExpectations(t)
}

func TestGetResourcesByHandler_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(mockResourceUsecase)
	mockNotify := new(mockNotifier)
	handler := rest.NewResourceHandler(mockUC, mockNotify)

	// Setup Gin
	r := gin.Default()
	r.GET("/customers/:id/resources", handler.GetResourcesByCustomer)

	mockUC.On("GetResourcesByCustomer", int64(1)).Return([]domain.Resource{
		{ID: 1, Name: "aws_vpc_main", Type: "VPC", Region: "us-east-1"}}, nil)

	req, _ := http.NewRequest("GET", "/customers/1/resources", nil)
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

	// Verify response content
	firstResource := data[0].(map[string]interface{})
	assert.Equal(t, "aws_vpc_main", firstResource["name"])
	assert.Equal(t, "VPC", firstResource["type"])
	assert.Equal(t, "us-east-1", firstResource["region"])

	mockUC.AssertExpectations(t)
}

func TestGetResourcesByHandler_InvalidCustomerID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mockResourceUsecase)
	mockNotify := new(mockNotifier)
	handler := rest.NewResourceHandler(mockUC, mockNotify)

	r := gin.Default()
	r.GET("/customers/:id/resources", handler.GetResourcesByCustomer)

	req, _ := http.NewRequest("GET", "/customers/abc/resources", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "invalid customer_id", resp["error"])

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetResourcesByHandler_CustomerNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(mockResourceUsecase)
	mockNotify := new(mockNotifier)
	handler := rest.NewResourceHandler(mockUC, mockNotify)

	// Setup Gin
	r := gin.Default()
	r.GET("/customers/:id/resources", handler.GetResourcesByCustomer)

	mockUC.On("GetResourcesByCustomer", int64(1)).Return(nil, errors.New("customer not found"))

	req, _ := http.NewRequest("GET", "/customers/1/resources", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "customer not found", resp["error"])

	mockUC.AssertExpectations(t)
}

func TestUpdateResourceHandler_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(mockResourceUsecase)
	mockNotify := new(mockNotifier)
	handler := rest.NewResourceHandler(mockUC, mockNotify)

	// Setup Gin
	r := gin.Default()
	r.PUT("/resources/:id", handler.UpdateResource)

	mockUC.On("UpdateResource", int64(1), "aws_vpc_main", "VPC", "us-east-1").Return(&domain.Resource{
		ID: 1, Name: "aws_vpc_main", Type: "VPC", Region: "us-east-1"}, nil)

	body := `{"name": "aws_vpc_main", "type": "VPC", "region": "us-east-1"}`
	req, _ := http.NewRequest("PUT", "/resources/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	mockUC.AssertExpectations(t)
}

func TestUpdateResourceHandler_InvalidResourceID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(mockResourceUsecase)
	mockNotify := new(mockNotifier)
	handler := rest.NewResourceHandler(mockUC, mockNotify)

	// Setup Gin
	r := gin.Default()
	r.PUT("/resources/:id", handler.UpdateResource)

	body := `{"name": "aws_vpc_main", "type": "VPC", "region": "us-east-1"}`
	req, _ := http.NewRequest("PUT", "/resources/abc", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "invalid resource id", resp["error"])

	mockUC.AssertExpectations(t)
}

func TestDeleteResourceHandler_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(mockResourceUsecase)
	mockNotify := new(mockNotifier)
	handler := rest.NewResourceHandler(mockUC, mockNotify)

	// Setup Gin
	r := gin.Default()
	r.DELETE("/resources/:id", handler.DeleteResource)

	mockUC.On("DeleteResource", int64(1)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/resources/1", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	mockUC.AssertExpectations(t)
}

func TestDeleteResourceHandler_InvalidResourceID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(mockResourceUsecase)
	mockNotify := new(mockNotifier)
	handler := rest.NewResourceHandler(mockUC, mockNotify)

	// Setup Gin
	r := gin.Default()
	r.DELETE("/resources/:id", handler.DeleteResource)

	req, _ := http.NewRequest("DELETE", "/resources/abc", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "invalid resource id", resp["error"])

	mockUC.AssertExpectations(t)
}
