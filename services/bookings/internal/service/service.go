package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

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
	uuidGenerator   func() uuid.UUID
}

func New(db database.Database,
	availabilitySvc availability.Availability,
	clock clockwork.Clock,
	uuidGenerator func() uuid.UUID) Service {
	return &service{
		db:              db,
		availabilitySvc: availabilitySvc,
		clock:           clock,
		uuidGenerator:   uuidGenerator,
	}
}

func (s *service) CreateBooking(ctx context.Context, create models.CreateBooking) (*models.Booking, error) {
	isAvailable, err := s.availabilitySvc.IsDateAvailable(ctx, create.LaunchPadID, create.LaunchDate)
	if err != nil {
		return nil, fmt.Errorf("cannot determine availability: %w", err)
	}
	if !isAvailable {
		return nil, models.ErrNotAvailable
	}
	now := s.clock.Now()
	result := models.Booking{
		ID:            s.uuidGenerator(),
		FirstName:     create.FirstName,
		LastName:      create.LastName,
		Gender:        create.Gender,
		Birthday:      create.Birthday,
		LaunchPadID:   create.LaunchPadID,
		DestinationID: create.DestinationID,
		LaunchDate:    create.LaunchDate,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	err = s.db.Create(ctx, result)
	if err != nil {
		return nil, fmt.Errorf("cannot create booking: %w", err)
	}
	return &result, nil
}

func (s *service) ListBookings(ctx context.Context, filters models.Filters, pagination models.Pagination) ([]models.Booking, error) {
	results, err := s.db.List(ctx, pagination, filters)
	if err != nil {
		return nil, fmt.Errorf("unable to list bookings: %w", err)
	}
	return results, nil
}

func (s *service) DeleteBooking(ctx context.Context, bookingID string) error {
	err := s.db.Delete(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("unable to delete booking: %w", err)
	}
	return nil
}
