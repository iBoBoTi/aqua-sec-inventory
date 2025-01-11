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

	"github.com/iBoBoTi/aqua-sec-inventory/internal/domain"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/transport/rest"
)

// Mock ResourceUsecase
type mockCustomerUsecase struct {
    mock.Mock
}

func (m *mockCustomerUsecase) CreateCustomer(name, email string)  (*domain.Customer, error) {
    args := m.Called(name, email)
    return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *mockCustomerUsecase) GetCustomerByID(customerID int64) (*domain.Customer, error) {
    args := m.Called(customerID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Customer), args.Error(1)
}

func TestCreateCustomerHandler_OK(t *testing.T) {
    gin.SetMode(gin.TestMode)

    mockUC := new(mockCustomerUsecase)
    handler := rest.NewCustomerHandler(mockUC)

    // Setup Gin
    r := gin.Default()
    r.POST("/customers", handler.CreateCustomer)

    mockUC.On("CreateCustomer", "ebuka", "test@email.com").Return(&domain.Customer{
		ID: 1,
		Name: "ebuka",
		Email: "test@email.com",
	},nil)

    body := `{"name":"ebuka","email":"test@email.com"}`
    req, _ := http.NewRequest("POST", "/customers", bytes.NewBufferString(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    // Perform request
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)

    var resp map[string]interface{}
    _ = json.Unmarshal(w.Body.Bytes(), &resp)

    // Verify response content
    customer := resp["data"].(map[string]interface{})
    assert.Equal(t, "ebuka", customer["name"])
    assert.Equal(t, "test@email.com", customer["email"])
    assert.NotZero(t,customer["id"])

    mockUC.AssertExpectations(t)
}

func TestAddCloudResourcesHandler_InternalServerError(t *testing.T) {
    gin.SetMode(gin.TestMode)

    mockUC := new(mockCustomerUsecase)
    handler := rest.NewCustomerHandler(mockUC)

    // Setup Gin
    r := gin.Default()
    r.POST("/customers", handler.CreateCustomer)

    mockUC.On("CreateCustomer", "ebuka", "test@email.com").Return((*domain.Customer)(nil),errors.New("internal server error"))

    body := `{"name":"ebuka","email":"test@email.com"}`
    req, _ := http.NewRequest("POST", "/customers", bytes.NewBufferString(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    // Perform request
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusInternalServerError, w.Code)

    var resp map[string]interface{}
    _ = json.Unmarshal(w.Body.Bytes(), &resp)

    mockUC.AssertExpectations(t)
}

func TestAddCloudResourcesHandler_EmptyName(t *testing.T) {
    gin.SetMode(gin.TestMode)

    mockUC := new(mockCustomerUsecase)
    handler := rest.NewCustomerHandler(mockUC)

    // Setup Gin
    r := gin.Default()
    r.POST("/customers", handler.CreateCustomer)

    mockUC.On("CreateCustomer", "ebuka", "test@email.com").Return((*domain.Customer)(nil),errors.New("name cannot be empty"))

    body := `{"name":"ebuka","email":"test@email.com"}`
    req, _ := http.NewRequest("POST", "/customers", bytes.NewBufferString(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    // Perform request
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

    var resp map[string]interface{}
    _ = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "name cannot be empty", resp["error"])

    mockUC.AssertExpectations(t)
}

func TestCreateCustomerHandler_EmptyEmail(t *testing.T) {
    gin.SetMode(gin.TestMode)

    mockUC := new(mockCustomerUsecase)
    handler := rest.NewCustomerHandler(mockUC)

    // Setup Gin
    r := gin.Default()
    r.POST("/customers", handler.CreateCustomer)

    mockUC.On("CreateCustomer", "ebuka", "test@email.com").Return((*domain.Customer)(nil),errors.New("email cannot be empty"))

    body := `{"name":"ebuka","email":"test@email.com"}`
    req, _ := http.NewRequest("POST", "/customers", bytes.NewBufferString(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    // Perform request
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

    var resp map[string]interface{}
    _ = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "email cannot be empty", resp["error"])

    mockUC.AssertExpectations(t)
}

func TestGetCustomerByIDHandler_OK(t *testing.T) {
    gin.SetMode(gin.TestMode)

    mockUC := new(mockCustomerUsecase)
    handler := rest.NewCustomerHandler(mockUC)

    // Setup Gin
    r := gin.Default()
    r.GET("/customers/:id", handler.GetCustomerByID)

    mockUC.On("GetCustomerByID", int64(1)).Return(&domain.Customer{
		ID: 1,
		Name: "ebuka",
		Email: "test@email.com",
	},nil)

    req, _ := http.NewRequest("GET", "/customers/1", nil)
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    // Perform request
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var resp map[string]interface{}
    _ = json.Unmarshal(w.Body.Bytes(), &resp)


    // Verify response content
    customer := resp["data"].(map[string]interface{})
    assert.Equal(t, "ebuka", customer["name"])
    assert.Equal(t, "test@email.com", customer["email"])
    assert.NotZero(t,customer["id"])

    mockUC.AssertExpectations(t)
}

func TestGetCustomerByIDHandler_InvalidCustomerID(t *testing.T) {
    gin.SetMode(gin.TestMode)

    mockUC := new(mockCustomerUsecase)
    handler := rest.NewCustomerHandler(mockUC)

    // Setup Gin
    r := gin.Default()
    r.GET("/customers/:id", handler.GetCustomerByID)

    req, _ := http.NewRequest("GET", "/customers/abc", nil)
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    // Perform request
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)

    var resp map[string]interface{}
    _ = json.Unmarshal(w.Body.Bytes(), &resp)

    assert.Equal(t, "invalid customer ID", resp["error"])

    mockUC.AssertExpectations(t)
}

func TestGetCustomerByIDHandler_CustomerNotFound(t *testing.T) {
    gin.SetMode(gin.TestMode)

    mockUC := new(mockCustomerUsecase)
    handler := rest.NewCustomerHandler(mockUC)

    // Setup Gin
    r := gin.Default()
    r.GET("/customers/:id", handler.GetCustomerByID)

    mockUC.On("GetCustomerByID", int64(1)).Return((*domain.Customer)(nil),errors.New("customer not found"))

    req, _ := http.NewRequest("GET", "/customers/1", nil)
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    // Perform request
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusNotFound, w.Code)

    var resp map[string]interface{}
    _ = json.Unmarshal(w.Body.Bytes(), &resp)


    // Verify response content
	
	assert.Equal(t, "customer not found", resp["error"])

    mockUC.AssertExpectations(t)
}