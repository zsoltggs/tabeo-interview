package database

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/models"
)

var ErrNotFound = errors.New("error not found")

//go:generate mockgen -package=mocks -destination=../mocks/database.go github.com/zsoltggs/tabeo-interview/services/bookings/internal/database Database
type Database interface {
	Create(ctx context.Context, booking models.Booking) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Booking, error)
	List(ctx context.Context, pagination models.Pagination, filters models.Filters) ([]models.Booking, error)
	Health() error
	Close(ctx context.Context)
}
