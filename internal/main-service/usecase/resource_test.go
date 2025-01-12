package usecase_test

import (
    "errors"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "github.com/iBoBoTi/aqua-sec-inventory/internal/main-service/domain"
    "github.com/iBoBoTi/aqua-sec-inventory/internal/main-service/usecase"
)

// Mock for ResourceRepository
type mockResourceRepo struct {
    mock.Mock
}

func (m *mockResourceRepo) AddResourcesToCustomer(resourceNames []string, customerID int64) error {
    args := m.Called(resourceNames, customerID)
    return args.Error(0)
}

func (m *mockResourceRepo) GetResourcesByCustomer(customerID int64) ([]domain.Resource, error) {
    args := m.Called(customerID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]domain.Resource), args.Error(1)
}

func (m *mockResourceRepo) GetByID(resourceID int64) (*domain.Resource, error) {
    args := m.Called(resourceID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Resource), args.Error(1)
}

func (m *mockResourceRepo) Update(r *domain.Resource) error {
    args := m.Called(r)
    return args.Error(0)
}

func (m *mockResourceRepo) Delete(resourceID int64) error {
    args := m.Called(resourceID)
    return args.Error(0)
}

func (m *mockResourceRepo) GetAll() ([]domain.Resource, error){
	args := m.Called()
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]domain.Resource), args.Error(1)
}

func (m *mockResourceRepo) GetByName(name string) (*domain.Resource, error) {
	args := m.Called(name)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Resource), args.Error(1)
}

func (m *mockResourceRepo) AddResourceToCustomer(resourceName string,customerID int64) error{
	args := m.Called(resourceName, customerID)
    return args.Error(0)
}

func (m *mockResourceRepo) GetCustomerResourceByResourceName(customerID int64, resourceName string) (*domain.Resource, error){

    args := m.Called(customerID, resourceName)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Resource), args.Error(1)
}

func (m *mockResourceRepo) DoesCustomerHaveResource(customerID int64, resourceName string) (bool, error) {
    args := m.Called(customerID, resourceName)
    if args.Get(0) == nil {
        return false, args.Error(1)
    }
    return args.Get(0).(bool), args.Error(1)
}

// Mock for CustomerRepository
type mockCustomerRepo2 struct {
    mock.Mock
}

func (m *mockCustomerRepo2) Create(customer *domain.Customer) error         { return nil }
func (m *mockCustomerRepo2) GetByID(id int64) (*domain.Customer, error)     {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Customer), args.Error(1)
}
func (m *mockCustomerRepo2) GetByEmail(email string) (*domain.Customer, error) {
    return nil, nil
}

func TestGetAllAvailableResourcesUsecase_OK(t *testing.T) {
    resourceRepo := new(mockResourceRepo)
    customerRepo := new(mockCustomerRepo2)

    uc := usecase.NewResourceUsecase(resourceRepo, customerRepo)

	resourceRepo.On("GetAll").Return([]domain.Resource{
		{ID: 1, Name: "aws_vpc_main", Type: "VPC", Region: "us-east-1"},
	},nil)

    resources,err := uc.GetAllAvailableResources()
    assert.NoError(t, err)
	assert.NotEmpty(t, resources)

    resourceRepo.AssertExpectations(t)
}

func TestAddCloudResourceUsecase_OK(t *testing.T) {
    resourceRepo := new(mockResourceRepo)
    customerRepo := new(mockCustomerRepo2)

    uc := usecase.NewResourceUsecase(resourceRepo, customerRepo)

    // Customer exists
    customerRepo.On("GetByID", int64(123)).Return(&domain.Customer{ID: 123}, nil)

    // Resource assignment
    resourceRepo.On("AddResourceToCustomer", "aws_vpc_main", int64(123)).
        Return(nil)

    resourceRepo.On("DoesCustomerHaveResource", int64(123),"aws_vpc_main").
        Return(false, nil)

    err := uc.AddCloudResource(123, "aws_vpc_main")
    assert.NoError(t, err)

    resourceRepo.AssertExpectations(t)
    customerRepo.AssertExpectations(t)
}

func TestAddCloudResourceUsecase_CustomerNotFound(t *testing.T) {
    resourceRepo := new(mockResourceRepo)
    customerRepo := new(mockCustomerRepo2)
    uc := usecase.NewResourceUsecase(resourceRepo, customerRepo)

    // Customer doesn't exist
    customerRepo.On("GetByID", int64(999)).
        Return((*domain.Customer)(nil), errors.New("no rows in result set"))

    err := uc.AddCloudResource(999, "resource1")
    assert.EqualError(t, err, "customer not found")

    customerRepo.AssertExpectations(t)
}

func TestAddCloudResourceUsecase_NoResourceNames (t *testing.T) {
    resourceRepo := new(mockResourceRepo)
    customerRepo := new(mockCustomerRepo2)
    uc := usecase.NewResourceUsecase(resourceRepo, customerRepo)

	customerRepo.On("GetByID", int64(123)).Return(&domain.Customer{ID: 123}, nil)

    err := uc.AddCloudResource(123, "")
    assert.EqualError(t, err, "no resource name provided")

    resourceRepo.AssertExpectations(t)
	customerRepo.AssertExpectations(t)
}

func TestGetResourcesByCustomerUsecase_OK(t *testing.T) {
    resourceRepo := new(mockResourceRepo)
    customerRepo := new(mockCustomerRepo2)

    uc := usecase.NewResourceUsecase(resourceRepo, customerRepo)

    // Customer exists
    customerRepo.On("GetByID", int64(123)).Return(&domain.Customer{ID: 123}, nil)

    // Resource assignment
    resourceRepo.On("GetResourcesByCustomer", int64(123)).
        Return([]domain.Resource{
			{ID: 1, Name: "aws_vpc_main", Type: "VPC", Region: "us-east-1"}},nil)

    resource,err := uc.GetResourcesByCustomer(123)
    assert.NoError(t, err)
	assert.NotEmpty(t, resource)

    resourceRepo.AssertExpectations(t)
    customerRepo.AssertExpectations(t)
}

func TestGetResourcesByCustomerUsecase_CustomerNotFound(t *testing.T) {
    resourceRepo := new(mockResourceRepo)
    customerRepo := new(mockCustomerRepo2)

    uc := usecase.NewResourceUsecase(resourceRepo, customerRepo)

    // Customer exists
	customerRepo.On("GetByID", int64(999)).
	Return((*domain.Customer)(nil), errors.New("no rows in result set"))


    _,err := uc.GetResourcesByCustomer(999)
	assert.EqualError(t, err, "customer not found")
    customerRepo.AssertExpectations(t)
}

func TestUpdateResourceUsecase_OK(t *testing.T) {
    resourceRepo := new(mockResourceRepo)
    customerRepo := new(mockCustomerRepo2)

    uc := usecase.NewResourceUsecase(resourceRepo, customerRepo)

	resourceRepo.On("GetByID",int64(1)).Return(&domain.Resource{
		ID: 1, Name: "aws_vpc_main", Type: "VPC", Region: "us-east-1",
	}, nil)

    // Resource assignment
    resourceRepo.On("Update", &domain.Resource{
		ID: 1, Name: "aws_vpc_main", Type: "VPC", Region: "us-east-1",
	}).Return(nil)

    
    _,err := uc.UpdateResource(1,"aws_vpc_main","VPC", "us-east-1")
    assert.NoError(t, err)

    resourceRepo.AssertExpectations(t)
    customerRepo.AssertExpectations(t)
}

func TestUpdateResourceUsecase_EmptyRegion(t *testing.T) {
    resourceRepo := new(mockResourceRepo)
    customerRepo := new(mockCustomerRepo2)

    uc := usecase.NewResourceUsecase(resourceRepo, customerRepo)

    _,err := uc.UpdateResource(1,"aws_vpc_main","VPC", "")
    assert.EqualError(t, err, "region cannot be empty")

}

func TestUpdateResourceUsecase_EmptyType(t *testing.T) {
    resourceRepo := new(mockResourceRepo)
    customerRepo := new(mockCustomerRepo2)

    uc := usecase.NewResourceUsecase(resourceRepo, customerRepo)

    _,err := uc.UpdateResource(1,"aws_vpc_main","", "us-east-1")
    assert.EqualError(t, err, "type cannot be empty")

}

func TestUpdateResourceUsecase_EmptyName(t *testing.T) {
    resourceRepo := new(mockResourceRepo)
    customerRepo := new(mockCustomerRepo2)

    uc := usecase.NewResourceUsecase(resourceRepo, customerRepo)
    
    _,err := uc.UpdateResource(1,"","VPC", "us-east-1")
    assert.EqualError(t, err, "name cannot be empty")

}

func TestUpdateResourceUsecase_ResourceNotFound(t *testing.T) {
    resourceRepo := new(mockResourceRepo)
    customerRepo := new(mockCustomerRepo2)

    uc := usecase.NewResourceUsecase(resourceRepo, customerRepo)

	resourceRepo.On("GetByID",int64(1)).Return((*domain.Resource)(nil), errors.New("no rows in result set"))

    _,err := uc.UpdateResource(1,"aws_vpc_main","VPC", "us-east-1")
    assert.EqualError(t, err, "resource not found")

    resourceRepo.AssertExpectations(t)
}

func TestDeleteResourceUsecase_OK(t *testing.T) {
    resourceRepo := new(mockResourceRepo)
    customerRepo := new(mockCustomerRepo2)

    uc := usecase.NewResourceUsecase(resourceRepo, customerRepo)

	resourceRepo.On("GetByID",int64(1)).Return(&domain.Resource{
		ID: 1, Name: "aws_vpc_main", Type: "VPC", Region: "us-east-1",
	}, nil)

	resourceRepo.On("Delete",int64(1)).Return(nil)
    
    err := uc.DeleteResource(1)
    assert.NoError(t, err)

    resourceRepo.AssertExpectations(t)
}

func TestDeleteResourceUsecase_ResourceNotFound(t *testing.T) {
    resourceRepo := new(mockResourceRepo)
    customerRepo := new(mockCustomerRepo2)

    uc := usecase.NewResourceUsecase(resourceRepo, customerRepo)

	resourceRepo.On("GetByID",int64(1)).Return((*domain.Resource)(nil), errors.New("no rows in result set"))

    err := uc.DeleteResource(1)
    assert.EqualError(t, err, "resource not found")

    resourceRepo.AssertExpectations(t)
}