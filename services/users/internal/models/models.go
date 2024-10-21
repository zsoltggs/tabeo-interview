package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Nickname  string `json:"nick_name"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Country   string `json:"country"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Filters struct {
	Email   *string `json:"email"`
	Country *string `json:"country"`
}

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
