package database

import (
	"context"
	"errors"

	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/models"
)

var ErrNotFound = errors.New("error not found")

//go:generate mockgen -package=mocks -destination=../mocks/database.go github.com/zsoltggs/tabeo-interview/services/bookings/internal/database Database
type Database interface {
	Create(ctx context.Context, user models.Booking) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*models.Booking, error)
	List(ctx context.Context, pagination models.Pagination, filters models.Filters) ([]models.Booking, error)
	Health() error
	Close(ctx context.Context)
}
