package customer

import "github.com/google/uuid"

type Customer struct {
	ID   int64     `json:"-"`
	UUID uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
