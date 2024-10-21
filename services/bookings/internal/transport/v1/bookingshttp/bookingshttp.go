package bookingshttp

import "net/http"

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
