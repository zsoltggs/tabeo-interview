// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package queries

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID            uuid.UUID
	FirstName     string
	LastName      string
	Gender        string
	Birthday      string
	LaunchPadID   string
	DestinationID string
	LaunchDate    time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
