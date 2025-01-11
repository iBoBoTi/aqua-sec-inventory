package repository

import (
    "database/sql"

    "github.com/iBoBoTi/aqua-sec-inventory/internal/domain"
)

type CustomerRepository interface {
    Create(customer *domain.Customer) error
    GetByID(id int64) (*domain.Customer, error)
    GetByEmail(email string) (*domain.Customer, error)
}

type customerRepo struct {
    db *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
    return &customerRepo{db: db}
}

func (r *customerRepo) Create(c *domain.Customer) error {
    query := `
        INSERT INTO customers (name, email, created_at, updated_at)
        VALUES ($1, $2, NOW(), NOW())
        RETURNING id
    `
    return r.db.QueryRow(query, c.Name, c.Email).Scan(&c.ID)
}

func (r *customerRepo) GetByID(id int64) (*domain.Customer, error) {
    query := `SELECT id, name, email, created_at, updated_at FROM customers WHERE id = $1`
    row := r.db.QueryRow(query, id)
    var c domain.Customer
    if err := row.Scan(&c.ID, &c.Name, &c.Email, &c.CreatedAt, &c.UpdatedAt); err != nil {
        return nil, err
    }
    return &c, nil
}

func (r *customerRepo) GetByEmail(email string) (*domain.Customer, error) {
    query := `SELECT id, name, email, created_at, updated_at FROM customers WHERE email = $1`
    row := r.db.QueryRow(query, email)
    var c domain.Customer
    if err := row.Scan(&c.ID, &c.Name, &c.Email, &c.CreatedAt, &c.UpdatedAt); err != nil {
        return nil, err
    }
    return &c, nil
}
