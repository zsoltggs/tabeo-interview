package database

import (
	"context"

	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/models"
)

type postgres struct {
}

func NewPostgres() (Database, error) {
	return &postgres{}, nil
}

func (p postgres) Create(ctx context.Context, user models.Booking) error {
	//TODO implement me
	panic("implement me")
}

func (p postgres) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (p postgres) GetByID(ctx context.Context, id string) (*models.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgres) List(ctx context.Context, pagination models.Pagination, filters models.Filters) ([]models.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgres) Health() error {
	//TODO implement me
	panic("implement me")
}

func (p postgres) Close(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}
