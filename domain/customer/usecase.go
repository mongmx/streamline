package customer

// UseCase - customer usecase APIs.
type UseCase interface {
	AllCustomer() ([]Customer, error)
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

func (u *useCase) CreateCustomer(customer *Customer) error {
	_, err := u.repo.CreateCustomer(customer)
	if err != nil {
		return err
	}
	return nil
}
