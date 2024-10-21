package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/models"
)

type pg struct {
	Pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, connectionStr string) (Database, error) {
	pool, err := pgxpool.New(ctx, connectionStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Check if the connection works
	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &pg{Pool: pool}, nil
}

func (p *pg) Create(ctx context.Context, user models.Booking) error {
	//TODO implement me
	panic("implement me")
}

func (p *pg) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (p pg) GetByID(ctx context.Context, id string) (*models.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (p pg) List(ctx context.Context, pagination models.Pagination, filters models.Filters) ([]models.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (p pg) Health() error {
	err := p.Pool.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}
	return nil
}

func (p pg) Close(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}
