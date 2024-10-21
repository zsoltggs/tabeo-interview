package bookingshttp

import (
	"net/http"

	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/models"
	bookingsv1 "github.com/zsoltggs/tabeo-interview/services/bookings/pkg/bookings/v1"
)

type BookingsHTTP interface {
	CreateBooking(response http.ResponseWriter, request *http.Request)
	ListBookings(response http.ResponseWriter, request *http.Request)
	DeleteBooking(response http.ResponseWriter, request *http.Request)
}

type service struct {
}

func New() BookingsHTTP {
	return &service{}
}

func (s service) CreateBooking(response http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s service) ListBookings(response http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s service) DeleteBooking(response http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}

func defaultPagination(pagination *bookingsv1.Pagination) models.Pagination {
	if pagination == nil {
		return models.Pagination{
			Offset: 0,
			Limit:  10,
		}
	}
	newPagination := models.Pagination{
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	}

	if newPagination.Limit == 0 {
		newPagination.Limit = 10
	}

	return newPagination
}
