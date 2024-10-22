package v1

import (
	"time"

	"github.com/google/uuid"
)

type CreateBookingRequest struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Gender        string `json:"gender"`
	Birthday      string `json:"birthday"`
	LaunchPadID   string `json:"launch_pad_id"`
	DestinationID string `json:"destination_id"`
	LaunchDate    string `json:"launch_date"`
}

type CreateBookingResponse struct {
	Booking *Booking `json:"booking,omitempty"`
	Error   string   `json:"error,omitempty"`
}

type Booking struct {
	ID uuid.UUID `json:"id"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Birthday  string `json:"birthday"`

	LaunchPadID   string `json:"launch_pad_id"`
	DestinationID string `json:"destination_id"`
	LaunchDate    string `json:"launch_date"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

type ListBookingsRequest struct {
	Filters    ListBookingsFilters `json:"filters"`
	Pagination Pagination          `json:"pagination"`
}

type ListBookingsResponse struct {
	Bookings []Booking `json:"bookings"`
}

type ListBookingsFilters struct {
	LaunchDate    *string `json:"launch_date"`
	LaunchPadID   *string `json:"launch_pad_id"`
	DestinationID *string `json:"destination_id"`
}

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
