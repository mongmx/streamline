package auth

import "github.com/google/uuid"

// UseCase - auth usecase APIs.
type UseCase interface {
	Register(user *User) error
	Signin(user *User) error
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

func (c useCase) Signin(user *User) error {
	sessionToken := uuid.New().String()
	err := c.repo.StoreSessionToken(sessionToken, user)
	if err != nil {
		return err
	}
	return nil
}
