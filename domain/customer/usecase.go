package customer

import (
	"context"
	"github.com/google/uuid"
)

// UseCase - customer usecase APIs.
type UseCase interface {
	AllCustomer() ([]Customer, error)
	GetCustomer(customerUUID uuid.UUID) (*Customer, error)
	StreamCustomer(ctx context.Context, id uuid.UUID, chatChan chan ChatMessage)
	CreateCustomer(customer *Customer) error
}

type useCase struct {
	repo Repository
}

// NewUseCase is a factory function of customer usecase.
func NewUseCase(repo Repository) UseCase {
	return &useCase{
		repo: repo,
	}
}

func (u *useCase) AllCustomer() ([]Customer, error) {
	return u.repo.AllCustomer()
}

func (u *useCase) GetCustomer(customerUUID uuid.UUID) (*Customer, error) {
	return u.repo.GetCustomer(customerUUID)
}

func (u *useCase) CreateCustomer(customer *Customer) error {
	_, err := u.repo.CreateCustomer(customer)
	if err != nil {
		return err
	}
	return nil
}

func (u *useCase) StreamCustomer(ctx context.Context, id uuid.UUID, chatChan chan ChatMessage) {
	u.repo.Streams(ctx, id, chatChan)
}
