package service

import (
	"context"

	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/service/availability"

	"github.com/jonboulle/clockwork"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/database"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/models"
)

//go:generate mockgen -package=mocks -destination=../mocks/service.go github.com/zsoltggs/tabeo-interview/services/bookings/internal/service Service
type Service interface {
	CreateBooking(ctx context.Context, createBooking models.CreateBooking) (*models.Booking, error)
	ListBookings(ctx context.Context, filters models.Filters, pagination models.Pagination) ([]models.Booking, error)
	DeleteBooking(ctx context.Context, bookingID string) error
}

type service struct {
	db              database.Database
	availabilitySvc availability.Availability
	clock           clockwork.Clock
}

func New(db database.Database,
	availabilitySvc availability.Availability,
	clock clockwork.Clock) Service {
	return &service{
		db:    db,
		clock: clock,
	}
}

func (s *service) CreateBooking(ctx context.Context, createBooking models.CreateBooking) (*models.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) ListBookings(ctx context.Context, filters models.Filters, pagination models.Pagination) ([]models.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) DeleteBooking(ctx context.Context, bookingID string) error {
	//TODO implement me
	panic("implement me")
}
