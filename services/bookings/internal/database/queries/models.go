// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package queries

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Booking struct {
	ID            uuid.UUID
	FirstName     string
	LastName      string
	Gender        string
	Birthday      pgtype.Timestamptz
	LaunchPadID   string
	DestinationID string
	LaunchDate    pgtype.Timestamptz
	CreatedAt     pgtype.Timestamptz
	UpdatedAt     pgtype.Timestamptz
}
