package bookingshttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/models"
	bookingsv1 "github.com/zsoltggs/tabeo-interview/services/bookings/pkg/bookings/v1"
)

func createListBookingsFromQueryParams(params url.Values) (*bookingsv1.ListBookingsRequest, error) {
	req := bookingsv1.ListBookingsRequest{}
	offset := 0
	offsetParam := params.Get("offset")
	if offsetParam != "" {
		parsedOffset, err := strconv.Atoi(offsetParam)
		if err != nil {
			return nil, errors.New("unable to parse offset")
		}
		offset = parsedOffset
		if offset < 0 {
			return nil, errors.New("invalid offset")
		}
	}

	req.Pagination.Offset = offset

	limit := 10 // Default limit
	limitParam := params.Get("limit")
	if limitParam != "" {
		parsedLimit, err := strconv.Atoi(params.Get("limit"))
		if err != nil {
			return nil, errors.New("unable to parse limit")
		}
		limit = parsedLimit
		if limit < 0 {
			return nil, errors.New("invalid limit")
		}
	}

	req.Pagination.Limit = limit

	if launchDateStr := params.Get("launch_date"); launchDateStr != "" {
		req.Filters.LaunchDate = &launchDateStr
		// TODO Ideally launch date should be in the future, but for now it will accept dates in the past to make
		// testing easier
	}

	if launchPadID := params.Get("launch_pad_id"); launchPadID != "" {
		req.Filters.LaunchPadID = &launchPadID
	}

	if destinationID := params.Get("destination_id"); destinationID != "" {
		req.Filters.DestinationID = &destinationID
	}
	return &req, nil
}

func toDomainFilter(filters bookingsv1.ListBookingsFilters) (models.Filters, error) {
	result := models.Filters{
		LaunchPadID:   filters.LaunchPadID,
		DestinationID: filters.DestinationID,
	}
	if filters.LaunchDate != nil {
		launchDate, err := time.Parse("2006-01-02", *filters.LaunchDate)
		if err != nil {
			return models.Filters{}, errors.New("invalid launch_date")
		}
		result.LaunchDate = &launchDate
	}
	return result, nil
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
	resp := bookingsv1.ErrorResponse{
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

// TODO Test
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
