package auth

import (
	"github.com/google/uuid"
	"time"
)

type model struct {
	ID        int64
	UUID      uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type User struct {
	model
	Email  string
	PlanID int64
	Auth   *Auth
	Topic  *Topic
}

type Auth struct {
	UserID int64
	Type   string
	Secret string
}

type Topic struct {
	model
	UserID int64
	Title  string
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"email"`
}
