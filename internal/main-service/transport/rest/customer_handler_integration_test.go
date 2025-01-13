package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/iBoBoTi/aqua-sec-inventory/internal/main-service/repository"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/main-service/transport/rest"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/main-service/usecase"
)

func TestCreateCustomerHandler_IntegrationTest_OK(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := setUpTestDB(t, "testdb", "testuser", "testpassword")
	defer db.Close()
	assert.NoError(t, err)

	r := gin.Default()
	repo := repository.NewCustomerRepository(db)
	resourceUC := usecase.NewCustomerUsecase(repo)
	handler := rest.NewCustomerHandler(resourceUC)

	r.POST("/customers", handler.CreateCustomer)

	// Perform the test request
	requestBody := map[string]string{"name": "testname", "email": "testemail@email.com"}
	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/customers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Validate response contents
	customer := response["data"].(map[string]interface{})

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "testname", customer["name"])
	assert.Equal(t, "testemail@email.com", customer["email"])
	assert.NotZero(t, customer["id"])

}

func TestCreateCustomerHandler_IntegrationTest_InvalidEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := setUpTestDB(t, "testdb", "testuser", "testpassword")
	defer db.Close()
	assert.NoError(t, err)

	r := gin.Default()
	repo := repository.NewCustomerRepository(db)
	resourceUC := usecase.NewCustomerUsecase(repo)
	handler := rest.NewCustomerHandler(resourceUC)

	r.POST("/customers", handler.CreateCustomer)

	// Perform the test request
	requestBody := map[string]string{"name": "testname", "email": "randomemail"}
	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/customers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestGetCustomerByIDHandler_IntegrationTest_OK(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := setUpTestDB(t, "testdb", "testuser", "testpassword")
	defer db.Close()
	assert.NoError(t, err)
	createdCustomer := seedCustomer(t, db)

	r := gin.Default()
	repo := repository.NewCustomerRepository(db)
	resourceUC := usecase.NewCustomerUsecase(repo)
	handler := rest.NewCustomerHandler(resourceUC)

	r.GET("/customers/:id", handler.GetCustomerByID)

	// Perform the test request
	url := fmt.Sprintf("/customers/%d", createdCustomer.ID)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	customer := response["data"].(map[string]interface{})

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, createdCustomer.Name, customer["name"])
	assert.Equal(t, createdCustomer.Email, customer["email"])

}

func TestGetCustomerByIDHandler_IntegrationTest_InvalidCustomerID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := setUpTestDB(t, "testdb", "testuser", "testpassword")
	defer db.Close()
	assert.NoError(t, err)

	r := gin.Default()
	repo := repository.NewCustomerRepository(db)
	resourceUC := usecase.NewCustomerUsecase(repo)
	handler := rest.NewCustomerHandler(resourceUC)

	r.GET("/customers/:id", handler.GetCustomerByID)

	// Perform the test request
	req, _ := http.NewRequest(http.MethodGet, "/customers/abc", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "invalid customer id", response["error"])

}

func TestGetCustomerByIDHandler_IntegrationTest_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := setUpTestDB(t, "testdb", "testuser", "testpassword")
	defer db.Close()
	assert.NoError(t, err)

	r := gin.Default()
	repo := repository.NewCustomerRepository(db)
	resourceUC := usecase.NewCustomerUsecase(repo)
	handler := rest.NewCustomerHandler(resourceUC)

	r.GET("/customers/:id", handler.GetCustomerByID)

	// Perform the test request
	url := fmt.Sprintf("/customers/%d", 1234)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "customer not found", response["error"])

}
