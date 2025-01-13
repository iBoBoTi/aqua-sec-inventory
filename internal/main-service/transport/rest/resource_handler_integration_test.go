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

	"github.com/iBoBoTi/aqua-sec-inventory/internal/main-service/domain"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/main-service/repository"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/main-service/transport/rest"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/main-service/usecase"
)

func TestGetAllResourceHandler_IntegrationTest_OK(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Create Customer Then Add a resource to the customer
	db, err := setUpTestDB(t, "testdb", "testuser", "testpassword")
	defer db.Close()
	assert.NoError(t, err)

	r := gin.Default()
	resourceRepo := repository.NewResourceRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	mockNotifer := new(mockNotifier)
	resourceUC := usecase.NewResourceUsecase(resourceRepo, customerRepo)
	handler := rest.NewResourceHandler(resourceUC, mockNotifer)

	resource1 := seedResource1(t, db)
	seedResource2(t, db)

	r.POST("/resources", handler.GetAllAvailableResources)

	// Perform the test request
	requestBody := map[string]string{"resource_name": resource1.Name}
	body, _ := json.Marshal(&requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/resources", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	resources := response["data"].([]interface{})

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 2, len(resources))

	mockNotifer.AssertExpectations(t)

}

func TestAddCloudResourceHandler_IntegrationTest_OK(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Create Customer Then Add a resource to the customer
	db, err := setUpTestDB(t, "testdb", "testuser", "testpassword")
	defer db.Close()
	assert.NoError(t, err)

	r := gin.Default()
	resourceRepo := repository.NewResourceRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	mockNotifer := new(mockNotifier)
	resourceUC := usecase.NewResourceUsecase(resourceRepo, customerRepo)
	handler := rest.NewResourceHandler(resourceUC, mockNotifer)

	cust := seedCustomer(t, db)
	resource1 := seedResource1(t, db)
	seedResource2(t, db)
	mockNotifer.On("Publish", domain.Notification{
		Event:   "notification",
		UserID:  cust.ID,
		Message: fmt.Sprintf("added resource %s for customer with customerID %d", resource1.Name, cust.ID),
	}).Return(nil)

	r.POST("/customers/:id/resources", handler.AddCloudResource)

	// Perform the test request
	requestBody := map[string]string{"resource_name": resource1.Name}
	body, _ := json.Marshal(&requestBody)

	url := fmt.Sprintf("/customers/%d/resources", cust.ID)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Validate response contents

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Resources assigned successfully", response["message"])

	ok, err := resourceRepo.DoesCustomerHaveResource(cust.ID, resource1.Name)
	assert.NoError(t, err)
	assert.True(t, ok)

	mockNotifer.AssertExpectations(t)

}

func TestAddCloudResourceHandler_IntegrationTest_ResourceAlreadyExist(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Create Customer Then Add a resource to the customer
	db, err := setUpTestDB(t, "testdb", "testuser", "testpassword")
	defer db.Close()
	assert.NoError(t, err)

	r := gin.Default()
	resourceRepo := repository.NewResourceRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	mockNotifer := new(mockNotifier)
	resourceUC := usecase.NewResourceUsecase(resourceRepo, customerRepo)
	handler := rest.NewResourceHandler(resourceUC, mockNotifer)

	cust := seedCustomer(t, db)
	resource1 := seedResource1(t, db)

	resourceRepo.AddResourceToCustomer(resource1.Name, cust.ID)

	r.POST("/customers/:id/resources", handler.AddCloudResource)

	// Perform the test request
	requestBody := map[string]string{"resource_name": resource1.Name}
	body, _ := json.Marshal(&requestBody)

	url := fmt.Sprintf("/customers/%d/resources", cust.ID)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Validate response contents

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, fmt.Sprintf("customer already has %s resource", resource1.Name), response["error"])

	mockNotifer.AssertExpectations(t)

}

func TestAddCloudResourceHandler_IntegrationTest_InvalidCustomerID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Create Customer Then Add a resource to the customer
	db, err := setUpTestDB(t, "testdb", "testuser", "testpassword")
	defer db.Close()
	assert.NoError(t, err)

	r := gin.Default()
	resourceRepo := repository.NewResourceRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	mockNotifer := new(mockNotifier)
	resourceUC := usecase.NewResourceUsecase(resourceRepo, customerRepo)
	handler := rest.NewResourceHandler(resourceUC, mockNotifer)

	resource1 := seedResource1(t, db)
	seedResource2(t, db)

	r.POST("/customers/:id/resources", handler.AddCloudResource)

	// Perform the test request
	requestBody := map[string]string{"resource_name": resource1.Name}
	body, _ := json.Marshal(&requestBody)

	url := fmt.Sprintf("/customers/%s/resources", "abc")
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Validate response contents

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "invalid customer id", response["error"])

	mockNotifer.AssertExpectations(t)

}

func TestGetResourcesByCustomerHandler_IntegrationTest_OK(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Create Customer Then Add a resource to the customer
	db, err := setUpTestDB(t, "testdb", "testuser", "testpassword")
	defer db.Close()
	assert.NoError(t, err)

	r := gin.Default()
	resourceRepo := repository.NewResourceRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	mockNotifer := new(mockNotifier)
	resourceUC := usecase.NewResourceUsecase(resourceRepo, customerRepo)
	handler := rest.NewResourceHandler(resourceUC, mockNotifer)

	cust := seedCustomer(t, db)
	resource1 := seedResource1(t, db)
	resource2 := seedResource2(t, db)
	resourceRepo.AddResourceToCustomer(resource1.Name, cust.ID)
	resourceRepo.AddResourceToCustomer(resource2.Name, cust.ID)

	r.GET("/customers/:id/resources", handler.GetResourcesByCustomer)

	// Perform the test request

	url := fmt.Sprintf("/customers/%d/resources", cust.ID)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	resources := response["data"].([]interface{})

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 2, len(resources))

	mockNotifer.AssertExpectations(t)

}

func TestGetResourcesByCustomerHandler_IntegrationTest_InvalidCustomerID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Create Customer Then Add a resource to the customer
	db, err := setUpTestDB(t, "testdb", "testuser", "testpassword")
	defer db.Close()
	assert.NoError(t, err)

	r := gin.Default()
	resourceRepo := repository.NewResourceRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	mockNotifer := new(mockNotifier)
	resourceUC := usecase.NewResourceUsecase(resourceRepo, customerRepo)
	handler := rest.NewResourceHandler(resourceUC, mockNotifer)

	r.GET("/customers/:id/resources", handler.GetResourcesByCustomer)

	// Perform the test request

	url := fmt.Sprintf("/customers/%s/resources", "abc")
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Validate response contents

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "invalid customer id", response["error"])

	mockNotifer.AssertExpectations(t)

}

func TestDeleteResourceHandler_IntegrationTest_OK(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Create Customer Then Add a resource to the customer
	db, err := setUpTestDB(t, "testdb", "testuser", "testpassword")
	defer db.Close()
	assert.NoError(t, err)

	r := gin.Default()
	resourceRepo := repository.NewResourceRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	mockNotifer := new(mockNotifier)
	resourceUC := usecase.NewResourceUsecase(resourceRepo, customerRepo)
	handler := rest.NewResourceHandler(resourceUC, mockNotifer)

	resource := seedResource1(t, db)
	r.DELETE("/resources/:id", handler.DeleteResource)

	// Perform the test request

	url := fmt.Sprintf("/resources/%d", resource.ID)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Validate response contents

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Resource deleted successfully", response["message"])

	res, err := resourceRepo.GetByID(resource.ID)
	assert.Error(t, err)
	assert.Nil(t, res)

	mockNotifer.AssertExpectations(t)

}

func TestDeleteResourceHandler_IntegrationTest_InvalidCustomerID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Create Customer Then Add a resource to the customer
	db, err := setUpTestDB(t, "testdb", "testuser", "testpassword")
	defer db.Close()
	assert.NoError(t, err)

	r := gin.Default()
	resourceRepo := repository.NewResourceRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	mockNotifer := new(mockNotifier)
	resourceUC := usecase.NewResourceUsecase(resourceRepo, customerRepo)
	handler := rest.NewResourceHandler(resourceUC, mockNotifer)

	r.DELETE("/resources/:id", handler.DeleteResource)

	// Perform the test request

	url := fmt.Sprintf("/resources/%s", "abc")
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Validate response contents

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "invalid resource id", response["error"])

	mockNotifer.AssertExpectations(t)

}

func TestUpdateResourceHandler_IntegrationTest_OK(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Create Customer Then Add a resource to the customer
	db, err := setUpTestDB(t, "testdb", "testuser", "testpassword")
	defer db.Close()
	assert.NoError(t, err)

	r := gin.Default()
	resourceRepo := repository.NewResourceRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	mockNotifer := new(mockNotifier)
	resourceUC := usecase.NewResourceUsecase(resourceRepo, customerRepo)
	handler := rest.NewResourceHandler(resourceUC, mockNotifer)

	resource := seedResource1(t, db)

	r.PUT("/resources/:id", handler.UpdateResource)

	// Perform the test request
	requestBody := map[string]string{"name": "new resource name", "type": "new type", "region": "new region"}
	body, _ := json.Marshal(&requestBody)

	url := fmt.Sprintf("/resources/%d", resource.ID)
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response domain.Resource
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Validate response contents

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "new resource name", response.Name)
	assert.Equal(t, "new type", response.Type)
	assert.Equal(t, "new region", response.Region)

	updatedResource, err := resourceRepo.GetByID(resource.ID)
	assert.NoError(t, err)

	assert.Equal(t, "new resource name", updatedResource.Name)
	assert.Equal(t, "new type", updatedResource.Type)
	assert.Equal(t, "new region", updatedResource.Region)

	mockNotifer.AssertExpectations(t)

}
func TestUpdateResourceHandler_IntegrationTest_InvalidCustomerID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Create Customer Then Add a resource to the customer
	db, err := setUpTestDB(t, "testdb", "testuser", "testpassword")
	defer db.Close()
	assert.NoError(t, err)

	r := gin.Default()
	resourceRepo := repository.NewResourceRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	mockNotifer := new(mockNotifier)
	resourceUC := usecase.NewResourceUsecase(resourceRepo, customerRepo)
	handler := rest.NewResourceHandler(resourceUC, mockNotifer)

	r.PUT("/resources/:id", handler.UpdateResource)

	// Perform the test request
	requestBody := map[string]string{"name": "new resource name", "type": "new type", "region": "new region"}
	body, _ := json.Marshal(&requestBody)

	url := fmt.Sprintf("/resources/%s", "abc")
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Validate response contents

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "invalid resource id", response["error"])

	mockNotifer.AssertExpectations(t)

}
