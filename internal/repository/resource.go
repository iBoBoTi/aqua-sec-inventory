package repository

import (
    "database/sql"
    "errors"

    "github.com/iBoBoTi/aqua-sec-inventory/internal/domain"
)

type ResourceRepository interface {
	GetAll() ([]domain.Resource, error)
    AddResourcesToCustomer(resourceNames []string, customerID int64) error
    GetResourcesByCustomer(customerID int64) ([]domain.Resource, error)
    GetByID(resourceID int64) (*domain.Resource, error)
    Update(resource *domain.Resource) error
    Delete(resourceID int64) error
    // Optionally: create or get resource by name
    GetByName(name string) (*domain.Resource, error)
}

type resourceRepo struct {
    db *sql.DB
}

func NewResourceRepository(db *sql.DB) ResourceRepository {
    return &resourceRepo{db: db}
}

func (r *resourceRepo) GetAll() ([]domain.Resource, error) {
	query := `SELECT id, name, type, region, customer_id, created_at, updated_at FROM resources`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []domain.Resource
	for rows.Next() {
		var res domain.Resource
		if err := rows.Scan(&res.ID, &res.Name, &res.Type, &res.Region, &res.CustomerID, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return nil, err
		}
		resources = append(resources, res)
	}
	return resources, nil
}

func (r *resourceRepo) GetByName(name string) (*domain.Resource, error) {
    query := `SELECT id, name, type, region, customer_id, created_at, updated_at
              FROM resources WHERE name = $1`
    row := r.db.QueryRow(query, name)
    var res domain.Resource
    if err := row.Scan(&res.ID, &res.Name, &res.Type, &res.Region, &res.CustomerID, &res.CreatedAt, &res.UpdatedAt); err != nil {
        return nil, err
    }
    return &res, nil
}

func (r *resourceRepo) AddResourcesToCustomer(resourceNames []string, customerID int64) error {
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }
    defer func() {
        if err != nil {
            _ = tx.Rollback()
        } else {
            _ = tx.Commit()
        }
    }()

    for _, name := range resourceNames {
        // Ensure resource exists
        resource, errGet := r.getResourceByNameTx(tx, name)
        if errGet != nil {
            return errors.New("resource " + name + " does not exist")
        }

        // Assign resource to customer
        updateQuery := `UPDATE resources SET customer_id = $1, updated_at = NOW() WHERE id = $2`
        _, errExec := tx.Exec(updateQuery, customerID, resource.ID)
        if errExec != nil {
            return errExec
        }
    }
    return nil
}

func (r *resourceRepo) getResourceByNameTx(tx *sql.Tx, name string) (*domain.Resource, error) {
    query := `SELECT id, name, type, region, customer_id, created_at, updated_at
              FROM resources WHERE name = $1`
    row := tx.QueryRow(query, name)
    var res domain.Resource
    if err := row.Scan(&res.ID, &res.Name, &res.Type, &res.Region, &res.CustomerID, &res.CreatedAt, &res.UpdatedAt); err != nil {
        return nil, err
    }
    return &res, nil
}

func (r *resourceRepo) GetResourcesByCustomer(customerID int64) ([]domain.Resource, error) {
    query := `SELECT id, name, type, region, customer_id, created_at, updated_at
              FROM resources
              WHERE customer_id = $1`
    rows, err := r.db.Query(query, customerID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var resources []domain.Resource
    for rows.Next() {
        var res domain.Resource
        err := rows.Scan(&res.ID, &res.Name, &res.Type, &res.Region, &res.CustomerID, &res.CreatedAt, &res.UpdatedAt)
        if err != nil {
            return nil, err
        }
        resources = append(resources, res)
    }
    return resources, nil
}

func (r *resourceRepo) GetByID(resourceID int64) (*domain.Resource, error) {
    query := `SELECT id, name, type, region, customer_id, created_at, updated_at
              FROM resources
              WHERE id = $1`
    row := r.db.QueryRow(query, resourceID)
    var res domain.Resource
    if err := row.Scan(&res.ID, &res.Name, &res.Type, &res.Region, &res.CustomerID, &res.CreatedAt, &res.UpdatedAt); err != nil {
        return nil, err
    }
    return &res, nil
}

func (r *resourceRepo) Update(resource *domain.Resource) error {
    query := `
        UPDATE resources
        SET name = $1, type = $2, region = $3, customer_id = $4, updated_at = NOW()
        WHERE id = $5
    `
    _, err := r.db.Exec(query, resource.Name, resource.Type, resource.Region, resource.CustomerID, resource.ID)
    return err
}

func (r *resourceRepo) Delete(resourceID int64) error {
    query := `DELETE FROM resources WHERE id = $1`
    _, err := r.db.Exec(query, resourceID)
    return err
}
