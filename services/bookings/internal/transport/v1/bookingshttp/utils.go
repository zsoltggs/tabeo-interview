package bookingshttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/models"
	bookingsv1 "github.com/zsoltggs/tabeo-interview/services/bookings/pkg/bookings/v1"
)

func defaultPagination(pagination *bookingsv1.Pagination) models.Pagination { //nolint
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

func fromDomainBooking(booking models.Booking) bookingsv1.Booking {
	return bookingsv1.Booking{
		ID:            booking.ID,
		FirstName:     booking.FirstName,
		LastName:      booking.LastName,
		Gender:        booking.Gender,
		Birthday:      booking.Birthday.Format("2006-01-02"),
		LaunchPadID:   booking.LaunchPadID,
		DestinationID: booking.DestinationID,
		LaunchDate:    booking.LaunchDate.Format("2006-01-02"),
		CreatedAt:     booking.CreatedAt,
		UpdatedAt:     booking.UpdatedAt,
	}
}

func writeErrorResponse(response http.ResponseWriter, reason string) {
	response.Header().Set("Content-Type", "application/json")
	resp := bookingsv1.CreateBookingResponse{
		Error: reason,
	}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		log.WithError(err).Error("unable to marshal create booking error response")
		return
	}
	_, err = response.Write(respJSON)
	if err != nil {
		log.WithError(err).Error("unable to write create booking error response")
		return
	}
}

func toDomainBooking(req bookingsv1.CreateBookingRequest) (*models.CreateBooking, error) {
	if req.DestinationID == "" {
		return nil, errors.New("destination id is required")
	}
	if req.LaunchPadID == "" {
		return nil, errors.New("launchpad id is required")
	}
	if req.Birthday == "" {
		return nil, errors.New("birthday is required")
	}
	if req.FirstName == "" {
		return nil, errors.New("first name is required")
	}
	if req.LastName == "" {
		return nil, errors.New("last name is required")
	}
	if req.Gender == "" {
		return nil, errors.New("gender is required")
	}
	if req.LaunchDate == "" {
		return nil, errors.New("launch date is required")
	}
	if req.Gender != "male" && req.Gender != "female" && req.Gender != "other" {
		return nil, errors.New("invalid gender value, accepted values for gender: male, female, other")
	}
	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		return nil, errors.New("invalid birthday, accepted format: 2006-01-02")
	}
	launchDate, err := time.Parse("2006-01-02", req.LaunchDate)
	if err != nil {
		return nil, errors.New("invalid birthday, accepted format: 2006-01-02")
	}
	result := models.CreateBooking{
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Gender:        req.Gender,
		Birthday:      birthday,
		LaunchPadID:   req.LaunchPadID,
		DestinationID: req.DestinationID,
		LaunchDate:    launchDate,
	}

	return &result, nil
}
