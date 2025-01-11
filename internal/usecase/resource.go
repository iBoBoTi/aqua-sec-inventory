package usecase

import (
	"errors"
	"fmt"
	"strings"

	"github.com/iBoBoTi/aqua-sec-inventory/internal/domain"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/repository"
)

type ResourceUsecase interface {
	GetAllAvailableResources() ([]domain.Resource, error)
    AddCloudResources(customerID int64, resourceNames []string) error
    GetResourcesByCustomer(customerID int64) ([]domain.Resource, error)
    UpdateResource(resourceID int64, name, resourceType, region string) (*domain.Resource, error)
    DeleteResource(resourceID int64) error
    AddCloudResource(customerID int64, resourceName string) error
}

type resourceUC struct {
    resourceRepo  repository.ResourceRepository
    customerRepo  repository.CustomerRepository
}

func NewResourceUsecase(resourceRepo repository.ResourceRepository, customerRepo repository.CustomerRepository) ResourceUsecase {
    return &resourceUC{
        resourceRepo: resourceRepo,
        customerRepo: customerRepo,
    }
}

func (uc *resourceUC) GetAllAvailableResources() ([]domain.Resource, error) {
	return uc.resourceRepo.GetAll()
}

func (uc *resourceUC) AddCloudResources(customerID int64, resourceNames []string) error {
    // Check if customer exists
    _, err := uc.customerRepo.GetByID(customerID)
    if err != nil {
        return errors.New("customer not found")
    }

    // Validate resourceNames
    if len(resourceNames) == 0 {
        return errors.New("no resource names provided")
    }

    return uc.resourceRepo.AddResourcesToCustomer(resourceNames, customerID)
}

func (uc *resourceUC) AddCloudResource(customerID int64, resourceName string) error {
    // Check if customer exists
    _, err := uc.customerRepo.GetByID(customerID)
    if err != nil {
        return errors.New("customer not found")
    }

    // Validate resourceNames
    if resourceName == "" {
        return errors.New("no resource name provided")
    }
    // Get customer resource by name if it exist send resource already exist error
    exist, err := uc.resourceRepo.DoesCustomerHaveResource(customerID ,resourceName)
    if err != nil {
        return err
    }

    if exist {
        return fmt.Errorf("customer already has %s resource", resourceName)
    }

    return uc.resourceRepo.AddResourceToCustomer(resourceName, customerID)
}

func (uc *resourceUC) GetResourcesByCustomer(customerID int64) ([]domain.Resource, error) {
    // Check if customer exists
    _, err := uc.customerRepo.GetByID(customerID)
    if err != nil {
        return nil, errors.New("customer not found")
    }

    return uc.resourceRepo.GetResourcesByCustomer(customerID)
}

func (uc *resourceUC) UpdateResource(resourceID int64, name, resourceType, region string) (*domain.Resource, error) {
    // Basic validations
    if strings.TrimSpace(name) == "" {
        return nil, errors.New("name cannot be empty")
    }
    if strings.TrimSpace(resourceType) == "" {
        return nil, errors.New("type cannot be empty")
    }
    if strings.TrimSpace(region) == "" {
        return nil, errors.New("region cannot be empty")
    }

    // Check if resource exists
    res, err := uc.resourceRepo.GetByID(resourceID)
    if err != nil {
        return nil, errors.New("resource not found")
    }

    // Update resource
    res.Name = name
    res.Type = resourceType
    res.Region = region

    if err := uc.resourceRepo.Update(res); err != nil {
        return nil, err
    }
    return res, nil
}

func (uc *resourceUC) DeleteResource(resourceID int64) error {
    // Check if resource exists
    _, err := uc.resourceRepo.GetByID(resourceID)
    if err != nil {
        return errors.New("resource not found")
    }

    return uc.resourceRepo.Delete(resourceID)
}
