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

func (p *useCase) Store(product *Product) error {
	return p.repo.Store(product)
}

func (p *useCase) Save(product *Product) error {
	return p.repo.Save(product)
}

func (p *useCase) StreamProduct(ctx context.Context, id uuid.UUID, prodChan chan Product) {
	p.repo.Streams(ctx, id, prodChan)
}
