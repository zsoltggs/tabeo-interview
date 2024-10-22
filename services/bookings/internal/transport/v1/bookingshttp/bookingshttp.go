package bookingshttp

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/service"

	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/models"
	bookingsv1 "github.com/zsoltggs/tabeo-interview/services/bookings/pkg/bookings/v1"
)

type BookingsHTTP interface {
	CreateBooking(response http.ResponseWriter, request *http.Request)
	ListBookings(response http.ResponseWriter, request *http.Request)
	DeleteBooking(response http.ResponseWriter, request *http.Request)
}

type bookingsHTTP struct {
	service service.Service
}

func New(service service.Service) BookingsHTTP {
	return &bookingsHTTP{
		service: service,
	}
}

func (h bookingsHTTP) CreateBooking(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(request.Body)
	if err != nil {
		writeErrorResponse(response, "bad request")
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	var bookingReq bookingsv1.CreateBookingRequest
	err = json.Unmarshal(body, &bookingReq)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		writeErrorResponse(response, "bad request")
		return
	}
	booking, err := toDomainBooking(bookingReq)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		writeErrorResponse(response, err.Error())
		return
	}
	ctx := request.Context()
	res, err := h.service.CreateBooking(ctx, *booking)
	switch {
	case errors.Is(err, models.ErrNotAvailable):
		response.WriteHeader(http.StatusConflict)
		writeErrorResponse(response, "date is unavailable")
		return
	case errors.Is(err, models.ErrNotFoundLaunchpad):
		response.WriteHeader(http.StatusNotFound)
		writeErrorResponse(response, "launch pad with ID not found")
		return
	case err != nil:
		response.WriteHeader(http.StatusInternalServerError)
		log.WithError(err).Error("unable to create booking")
		return
	}
	result := fromDomainBooking(*res)
	resp := bookingsv1.CreateBookingResponse{
		Booking: &result,
	}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		log.WithError(err).Error("unable to marshal create booking response")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	_, err = response.Write(respJSON)
	if err != nil {
		log.WithError(err).Error("unable to write create booking response")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func (h bookingsHTTP) ListBookings(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req, err := createListBookingsFromQueryParams(request.URL.Query())
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		writeErrorResponse(response, err.Error())
		return
	}

	filters, err := toDomainFilter(req.Filters)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		writeErrorResponse(response, err.Error())
		return
	}

	ctx := request.Context()
	bookings, err := h.service.ListBookings(ctx, filters, models.Pagination{
		Offset: req.Pagination.Offset,
		Limit:  req.Pagination.Limit,
	})
	if err != nil {
		log.WithError(err).Error("unable to list bookings")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	var results []bookingsv1.Booking
	for _, b := range bookings {
		results = append(results, fromDomainBooking(b))
	}

	resp := bookingsv1.ListBookingsResponse{
		Bookings: results,
	}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		log.WithError(err).Error("unable to marshal list bookings response")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	_, err = response.Write(respJSON)
	if err != nil {
		log.WithError(err).Error("unable to write list bookings response")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h bookingsHTTP) DeleteBooking(response http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}
