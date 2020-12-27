package product

import (
	"context"

	"github.com/google/uuid"
)

// UseCase - product usecase APIs.
type UseCase interface {
	Store(product *Product) error
	Save(product *Product) error
	StreamProduct(ctx context.Context, id uuid.UUID, prodChan chan Product)
}

type useCase struct {
	repo Repository
}

// NewUseCase is a factory function of product usecase.
func NewUseCase(repo Repository) UseCase {
	return &useCase{
		repo: repo,
	}
}

func (u *useCase) Store(product *Product) error {
	return u.repo.Store(product)
}

func (u *useCase) Save(product *Product) error {
	return u.repo.Save(product)
}

func (u *useCase) StreamProduct(ctx context.Context, id uuid.UUID, prodChan chan Product) {
	u.repo.Streams(ctx, id, prodChan)
}
