package models

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID uuid.UUID `json:"id"`

	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Gender    string    `json:"gender"`
	Birthday  time.Time `json:"birthday"`

	LaunchPadID   string    `json:"launch_pad_id"`
	DestinationID string    `json:"destination_id"`
	LaunchDate    time.Time `json:"launch_date"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateBooking struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Gender    string    `json:"gender"`
	BirthDay  time.Time `json:"birthday"`

	LaunchPadID   string    `json:"launch_pad_id"`
	DestinationID string    `json:"destination_id"`
	LaunchDate    time.Time `json:"launch_date"`
}

type Filters struct {
	Email   *string `json:"email"`
	Country *string `json:"country"`
}

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
