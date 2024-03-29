package auth

// UseCase - auth usecase APIs.
type UseCase interface {
	Register(user *User) error
	Signin(sessionToken string, user *User) error
	Profile(sessionToken string) (User, error)
}

type useCase struct {
	repo Repository
}

// NewUseCase is a factory function of auth usecase.
func NewUseCase(repo Repository) UseCase {
	return &useCase{
		repo: repo,
	}
}

func (c useCase) Register(user *User) error {
	user, err := c.repo.CreateUser(user)
	if err != nil {
		return err
	}
	user.Auth.UserID = user.ID
	_, err = c.repo.CreateAuth(user.Auth)
	if err != nil {
		return err
	}
	user.Topic.UserID = user.ID
	_, err = c.repo.CreateTopic(user.Topic)
	if err != nil {
		return err
	}
	return nil
}

func (c useCase) Signin(sessionToken string, user *User) error {
	err := c.repo.StoreSessionUser(sessionToken, user)
	if err != nil {
		return err
	}
	return nil
}

func (c useCase) Profile(sessionToken string) (User, error) {
	user, err := c.repo.GetSessionUser(sessionToken)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
