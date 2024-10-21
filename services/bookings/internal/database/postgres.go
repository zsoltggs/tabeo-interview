package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/database/queries"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/models"
)

type pg struct {
	pool    *pgxpool.Pool
	queries *queries.Queries
}

func NewPostgres(ctx context.Context, connectionStr string) (Database, error) {
	pool, err := pgxpool.New(ctx, connectionStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &pg{
		pool:    pool,
		queries: queries.New(pool),
	}, nil
}

func (q *pg) Create(ctx context.Context, booking models.Booking) error {
	err := q.queries.CreateBooking(ctx, queries.CreateBookingParams{
		ID:            booking.ID,
		FirstName:     booking.FirstName,
		LastName:      booking.LastName,
		Gender:        booking.Gender,
		Birthday:      pgtype.Timestamptz{Time: booking.Birthday, Valid: true},
		LaunchPadID:   booking.LaunchPadID,
		DestinationID: booking.DestinationID,
		LaunchDate:    pgtype.Timestamptz{Time: booking.LaunchDate, Valid: true},
		CreatedAt:     pgtype.Timestamptz{Time: booking.CreatedAt, Valid: true},
		UpdatedAt:     pgtype.Timestamptz{Time: booking.UpdatedAt, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("error creating booking: %w", err)
	}
	return nil
}

func (q *pg) Delete(ctx context.Context, id uuid.UUID) error {
	err := q.queries.DeleteBooking(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting booking: %w", err)
	}
	return nil
}

func (q *pg) GetByID(ctx context.Context, id uuid.UUID) (*models.Booking, error) {
	booking, err := q.queries.GetBookingByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &models.Booking{
		ID:            booking.ID,
		FirstName:     booking.FirstName,
		LastName:      booking.LastName,
		Gender:        booking.Gender,
		Birthday:      booking.Birthday.Time, // All values are required so this is fine
		LaunchPadID:   booking.LaunchPadID,
		DestinationID: booking.DestinationID,
		LaunchDate:    booking.LaunchDate.Time,
		CreatedAt:     booking.CreatedAt.Time,
		UpdatedAt:     booking.UpdatedAt.Time,
	}, nil
}

func (q *pg) List(ctx context.Context, pagination models.Pagination, filters models.Filters) ([]models.Booking, error) {
	params := queries.ListBookingsParams{
		LaunchDate:    pgtype.Timestamptz{},
		LaunchPadID:   pgtype.Text{},
		DestinationID: pgtype.Text{},
		Offset:        int32(pagination.Offset),
		Limit:         int32(pagination.Limit),
	}
	if filters.LaunchDate != nil {
		params.LaunchDate = pgtype.Timestamptz{
			Time:  *filters.LaunchDate,
			Valid: true,
		}
	}
	if filters.DestinationID != nil {
		params.DestinationID = pgtype.Text{
			String: *filters.DestinationID,
			Valid:  true,
		}
	}
	if filters.LaunchPadID != nil {
		params.LaunchPadID = pgtype.Text{
			String: *filters.LaunchPadID,
			Valid:  true,
		}
	}
	bookings, err := q.queries.ListBookings(ctx, params)
	if err != nil {
		return nil, err
	}
	var result []models.Booking
	for _, b := range bookings {
		result = append(result, models.Booking{
			ID:            b.ID,
			FirstName:     b.FirstName,
			LastName:      b.LastName,
			Gender:        b.Gender,
			Birthday:      b.Birthday.Time, // All values are required so this is fine
			LaunchPadID:   b.LaunchPadID,
			DestinationID: b.DestinationID,
			LaunchDate:    b.LaunchDate.Time,
			CreatedAt:     b.CreatedAt.Time,
			UpdatedAt:     b.UpdatedAt.Time,
		})
	}
	return result, nil
}

func (q *pg) Health() error {
	return q.pool.Ping(context.Background())
}

func (q *pg) Close(ctx context.Context) {
	q.pool.Close()
}
