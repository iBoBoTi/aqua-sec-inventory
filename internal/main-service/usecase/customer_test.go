package usecase_test

import (
	"errors"
	"testing"

	"github.com/iBoBoTi/aqua-sec-inventory/internal/main-service/domain"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/main-service/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository
type mockCustomerRepo struct {
	mock.Mock
}

func (m *mockCustomerRepo) Create(customer *domain.Customer) error {
	args := m.Called(customer)
	return args.Error(0)
}

func (m *mockCustomerRepo) GetByID(id int64) (*domain.Customer, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *mockCustomerRepo) GetByEmail(email string) (*domain.Customer, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func TestCreateCustomer_OK(t *testing.T) {
	repo := new(mockCustomerRepo)
	uc := usecase.NewCustomerUsecase(repo)

	repo.On("GetByEmail", "john@example.com").Return((*domain.Customer)(nil), errors.New("not found"))
	repo.On("Create", mock.AnythingOfType("*domain.Customer")).Return(nil)

	cust, err := uc.CreateCustomer("John", "john@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, cust)
	assert.Equal(t, "John", cust.Name)
	assert.Equal(t, "john@example.com", cust.Email)

	repo.AssertExpectations(t)
}

func TestCreateCustomer_DuplicateEmail(t *testing.T) {
	repo := new(mockCustomerRepo)
	uc := usecase.NewCustomerUsecase(repo)

	existingCust := &domain.Customer{ID: 1, Name: "Existing", Email: "john@example.com"}
	repo.On("GetByEmail", "john@example.com").Return(existingCust, nil)

	cust, err := uc.CreateCustomer("John", "john@example.com")
	assert.Nil(t, cust)
	assert.EqualError(t, err, "customer with this email already exists")
}

func TestCreateCustomer_EmptyName(t *testing.T) {
	repo := new(mockCustomerRepo)
	uc := usecase.NewCustomerUsecase(repo)

	cust, err := uc.CreateCustomer("", "john@example.com")
	assert.EqualError(t, err, "name cannot be empty")
	assert.Nil(t, cust)
}

func TestCreateCustomer_EmptyEmail(t *testing.T) {
	repo := new(mockCustomerRepo)
	uc := usecase.NewCustomerUsecase(repo)

	cust, err := uc.CreateCustomer("John", "")
	assert.EqualError(t, err, "email cannot be empty")
	assert.Nil(t, cust)

}

func TestGetCustomerByIDUsecase_OK(t *testing.T) {
	customerRepo := new(mockCustomerRepo2)

	uc := usecase.NewCustomerUsecase(customerRepo)

	// Customer exists
	customerRepo.On("GetByID", int64(123)).Return(&domain.Customer{ID: 123}, nil)

	resource, err := uc.GetCustomerByID(123)
	assert.NoError(t, err)
	assert.NotEmpty(t, resource)

	customerRepo.AssertExpectations(t)
}
