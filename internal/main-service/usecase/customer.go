package usecase

import (
	"errors"
	"log"
	"strings"

	"github.com/iBoBoTi/aqua-sec-inventory/internal/main-service/domain"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/main-service/repository"
)

type CustomerUsecase interface {
    CreateCustomer(name, email string) (*domain.Customer, error)
    GetCustomerByID(id int64) (*domain.Customer, error)
}

type customerUC struct {
    customerRepo repository.CustomerRepository
}

func NewCustomerUsecase(customerRepo repository.CustomerRepository) CustomerUsecase {
    return &customerUC{
        customerRepo: customerRepo,
    }
}

func (uc *customerUC) CreateCustomer(name, email string) (*domain.Customer, error) {
    // Basic validation
    if strings.TrimSpace(name) == "" {
        return nil, errors.New("name cannot be empty")
    }
    if strings.TrimSpace(email) == "" {
        return nil, errors.New("email cannot be empty")
    }

    // Check if email already exists
    existing, _ := uc.customerRepo.GetByEmail(email)
    if existing != nil {
        return nil, errors.New("customer with this email already exists")
    }

    c := &domain.Customer{
        Name:  name,
        Email: email,
    }
    if err := uc.customerRepo.Create(c); err != nil {
		log.Println("Error creating customer: ", err)
        return nil, errors.New("internal server error")
    }
	
    return c, nil
}

func (uc *customerUC) GetCustomerByID(id int64) (*domain.Customer, error) {
    return uc.customerRepo.GetByID(id)
}
