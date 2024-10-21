package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/mocks"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/models"
	"go.uber.org/mock/gomock"
)

func TestService_CreateBooking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	mockAvailabilitySvc := mocks.NewMockAvailability(ctrl)
	mockedTime := time.Date(2024, 01, 01, 01, 1, 1, 1, time.UTC)
	ts := time.Now().Truncate(time.Second)
	mockClock := clockwork.NewFakeClockAt(mockedTime)
	fixedID := "65383d1f-ef0f-4250-893b-4c72c91f4b25"
	mockUUID := uuid.MustParse(fixedID)
	uuidGen := func() uuid.UUID { return mockUUID }
	const validLunchPadID = "5e9e4501f509094ba4566f84"
	expectedValidBooking := models.Booking{
		ID:            mockUUID,
		FirstName:     "John",
		LastName:      "Doe",
		Gender:        "male",
		Birthday:      ts,
		LaunchPadID:   validLunchPadID,
		DestinationID: "destination_1",
		LaunchDate:    ts,
		CreatedAt:     mockedTime,
		UpdatedAt:     mockedTime,
	}

	svc := New(mockDB, mockAvailabilitySvc, mockClock, uuidGen)

	tests := []struct {
		name            string
		input           models.CreateBooking
		mockSetup       func()
		expectedBooking *models.Booking
		expectedError   error
	}{
		{
			name: "Successful booking creation",
			input: models.CreateBooking{
				FirstName:     "John",
				LastName:      "Doe",
				Gender:        "male",
				Birthday:      ts,
				LaunchPadID:   validLunchPadID,
				DestinationID: "destination_1",
				LaunchDate:    ts,
			},
			mockSetup: func() {
				mockAvailabilitySvc.EXPECT().
					IsDateAvailable(gomock.Any(), validLunchPadID, ts).
					Return(true, nil)
				mockDB.EXPECT().
					Create(gomock.Any(), expectedValidBooking).
					Return(nil)
			},
			expectedBooking: &expectedValidBooking,
			expectedError:   nil,
		},
		{
			name: "Date not available",
			input: models.CreateBooking{
				LaunchPadID: validLunchPadID,
				LaunchDate:  ts,
			},
			mockSetup: func() {
				mockAvailabilitySvc.EXPECT().
					IsDateAvailable(gomock.Any(), validLunchPadID, ts).
					Return(false, nil)
			},
			expectedBooking: nil,
			expectedError:   models.ErrNotAvailable,
		},
		{
			name: "Availability service returns an error",
			input: models.CreateBooking{
				LaunchPadID: validLunchPadID,
				LaunchDate:  ts,
			},
			mockSetup: func() {
				mockAvailabilitySvc.EXPECT().
					IsDateAvailable(gomock.Any(), validLunchPadID, ts).
					Return(false, errors.New("service unavailable"))
			},
			expectedBooking: nil,
			expectedError:   errors.New("cannot determine availability: service unavailable"),
		},
		{
			name: "Error creating booking in database",
			input: models.CreateBooking{
				FirstName:     "John",
				LastName:      "Doe",
				Gender:        "male",
				Birthday:      time.Now(),
				LaunchPadID:   validLunchPadID,
				DestinationID: "destination_1",
				LaunchDate:    ts,
			},
			mockSetup: func() {
				mockAvailabilitySvc.EXPECT().
					IsDateAvailable(gomock.Any(), validLunchPadID, ts).
					Return(true, nil)
				mockDB.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errors.New("database error"))
			},
			expectedBooking: nil,
			expectedError:   errors.New("cannot create booking: database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			booking, err := svc.CreateBooking(context.Background(), tt.input)

			assert.Equal(t, tt.expectedBooking, booking)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_ListBookings(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	mockAvailabilitySvc := mocks.NewMockAvailability(ctrl)
	mockClock := clockwork.NewRealClock()
	uuidGen := func() uuid.UUID { return uuid.New() }

	svc := New(mockDB, mockAvailabilitySvc, mockClock, uuidGen)

	tests := []struct {
		name          string
		filters       models.Filters
		pagination    models.Pagination
		mockSetup     func()
		expectedList  []models.Booking
		expectedError error
	}{
		{
			name: "List bookings successfully",
			filters: models.Filters{
				LaunchDate:    toPtr("a"),
				LaunchPadID:   toPtr("b"),
				DestinationID: toPtr("c"),
			},
			pagination: models.Pagination{Limit: 10, Offset: 1},
			mockSetup: func() {
				mockDB.EXPECT().
					List(gomock.Any(),
						models.Pagination{Limit: 10, Offset: 1}, models.Filters{
							LaunchDate:    toPtr("a"),
							LaunchPadID:   toPtr("b"),
							DestinationID: toPtr("c"),
						}).
					Return([]models.Booking{{FirstName: "John"}}, nil)
			},
			expectedList:  []models.Booking{{FirstName: "John"}},
			expectedError: nil,
		},
		{
			name:       "Error listing bookings",
			filters:    models.Filters{},
			pagination: models.Pagination{Limit: 10, Offset: 1},
			mockSetup: func() {
				mockDB.EXPECT().
					List(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("list error"))
			},
			expectedList:  nil,
			expectedError: errors.New("unable to list bookings: list error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			bookings, err := svc.ListBookings(context.Background(), tt.filters, tt.pagination)

			assert.Equal(t, tt.expectedList, bookings)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_DeleteBooking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	mockAvailabilitySvc := mocks.NewMockAvailability(ctrl)
	mockClock := clockwork.NewFakeClock()
	uuidGen := func() uuid.UUID { return uuid.New() }

	svc := New(mockDB, mockAvailabilitySvc, mockClock, uuidGen)

	tests := []struct {
		name          string
		bookingID     string
		mockSetup     func()
		expectedError error
	}{
		{
			name:      "Successful delete",
			bookingID: "booking_1",
			mockSetup: func() {
				mockDB.EXPECT().
					Delete(gomock.Any(), "booking_1").
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "Error deleting booking",
			bookingID: "booking_1",
			mockSetup: func() {
				mockDB.EXPECT().
					Delete(gomock.Any(), "booking_1").
					Return(errors.New("delete error"))
			},
			expectedError: errors.New("unable to delete booking: delete error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := svc.DeleteBooking(context.Background(), tt.bookingID)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func toPtr(s string) *string {
	return &s
}
