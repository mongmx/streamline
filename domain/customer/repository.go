package customer

import (
	"database/sql"
	"github.com/google/uuid"
)

// Repository - customer store APIs.
type Repository interface {
	AllCustomer() ([]Customer, error)
	CreateCustomer(customer *Customer) (int64, error)
}

type repo struct {
	DB *sql.DB
}

// NewRepository is a factory function of customer store.
func NewRepository(db *sql.DB) Repository {
	return &repo{
		DB: db,
	}
}

func (r *repo) AllCustomer() ([]Customer, error) {
	customers := make([]Customer, 0)
	customer := Customer{
		ID:   1,
		UUID: uuid.New(),
		Name: "John Doe",
	}
	customers = append(customers, customer)
	return customers, nil
}

func (r *repo) CreateCustomer(customer *Customer) (int64, error) {
	return customer.ID, nil
}
