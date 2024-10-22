package bookingshttp

import (
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/models"
	bookingsv1 "github.com/zsoltggs/tabeo-interview/services/bookings/pkg/bookings/v1"
)

func TestToDomainBooking(t *testing.T) {
	tests := []struct {
		name        string
		input       bookingsv1.CreateBookingRequest
		expected    *models.CreateBooking
		expectedErr error
	}{
		{
			name: "valid input",
			input: bookingsv1.CreateBookingRequest{
				FirstName:     "John",
				LastName:      "Doe",
				Gender:        "male",
				Birthday:      "1990-01-01",
				LaunchPadID:   "lp-123",
				DestinationID: "dest-456",
				LaunchDate:    "2024-12-01",
			},
			expected: &models.CreateBooking{
				FirstName:     "John",
				LastName:      "Doe",
				Gender:        "male",
				Birthday:      parseDate("1990-01-01"),
				LaunchPadID:   "lp-123",
				DestinationID: "dest-456",
				LaunchDate:    parseDate("2024-12-01"),
			},
			expectedErr: nil,
		},
		{
			name: "missing destination id",
			input: bookingsv1.CreateBookingRequest{
				FirstName:   "John",
				LastName:    "Doe",
				Gender:      "male",
				Birthday:    "1990-01-01",
				LaunchPadID: "lp-123",
				LaunchDate:  "2024-12-01",
			},
			expected:    nil,
			expectedErr: errors.New("destination id is required"),
		},
		{
			name: "missing launchpad id",
			input: bookingsv1.CreateBookingRequest{
				FirstName:     "John",
				LastName:      "Doe",
				Gender:        "male",
				Birthday:      "1990-01-01",
				DestinationID: "dest-456",
				LaunchDate:    "2024-12-01",
			},
			expected:    nil,
			expectedErr: errors.New("launchpad id is required"),
		},
		{
			name: "missing birthday",
			input: bookingsv1.CreateBookingRequest{
				FirstName:     "John",
				LastName:      "Doe",
				Gender:        "male",
				LaunchPadID:   "lp-123",
				DestinationID: "dest-456",
				LaunchDate:    "2024-12-01",
			},
			expected:    nil,
			expectedErr: errors.New("birthday is required"),
		},
		{
			name: "missing first name",
			input: bookingsv1.CreateBookingRequest{
				LastName:      "Doe",
				Gender:        "male",
				Birthday:      "1990-01-01",
				LaunchPadID:   "lp-123",
				DestinationID: "dest-456",
				LaunchDate:    "2024-12-01",
			},
			expected:    nil,
			expectedErr: errors.New("first name is required"),
		},
		{
			name: "missing last name",
			input: bookingsv1.CreateBookingRequest{
				FirstName:     "John",
				Gender:        "male",
				Birthday:      "1990-01-01",
				LaunchPadID:   "lp-123",
				DestinationID: "dest-456",
				LaunchDate:    "2024-12-01",
			},
			expected:    nil,
			expectedErr: errors.New("last name is required"),
		},
		{
			name: "missing gender",
			input: bookingsv1.CreateBookingRequest{
				FirstName:     "John",
				LastName:      "Doe",
				Birthday:      "1990-01-01",
				LaunchPadID:   "lp-123",
				DestinationID: "dest-456",
				LaunchDate:    "2024-12-01",
			},
			expected:    nil,
			expectedErr: errors.New("gender is required"),
		},
		{
			name: "missing launch date",
			input: bookingsv1.CreateBookingRequest{
				FirstName:     "John",
				LastName:      "Doe",
				Gender:        "male",
				Birthday:      "1990-01-01",
				LaunchPadID:   "lp-123",
				DestinationID: "dest-456",
			},
			expected:    nil,
			expectedErr: errors.New("launch date is required"),
		},
		{
			name: "invalid gender value",
			input: bookingsv1.CreateBookingRequest{
				FirstName:     "John",
				LastName:      "Doe",
				Gender:        "invalid",
				Birthday:      "1990-01-01",
				LaunchPadID:   "lp-123",
				DestinationID: "dest-456",
				LaunchDate:    "2024-12-01",
			},
			expected:    nil,
			expectedErr: errors.New("invalid gender value, accepted values for gender: male, female, other"),
		},
		{
			name: "invalid birthday format",
			input: bookingsv1.CreateBookingRequest{
				FirstName:     "John",
				LastName:      "Doe",
				Gender:        "male",
				Birthday:      "invalid-date",
				LaunchPadID:   "lp-123",
				DestinationID: "dest-456",
				LaunchDate:    "2024-12-01",
			},
			expected:    nil,
			expectedErr: errors.New("invalid birthday, accepted format: 2006-01-02"),
		},
		{
			name: "invalid launch date format",
			input: bookingsv1.CreateBookingRequest{
				FirstName:     "John",
				LastName:      "Doe",
				Gender:        "male",
				Birthday:      "1990-01-01",
				LaunchPadID:   "lp-123",
				DestinationID: "dest-456",
				LaunchDate:    "invalid-date",
			},
			expected:    nil,
			expectedErr: errors.New("invalid birthday, accepted format: 2006-01-02"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := toDomainBooking(tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestCreateListBookingsFromQueryParams(t *testing.T) {
	tests := []struct {
		name        string
		params      url.Values
		expected    *bookingsv1.ListBookingsRequest
		expectedErr error
	}{
		{
			name: "valid parameters",
			params: url.Values{
				"offset":         []string{"5"},
				"limit":          []string{"20"},
				"launch_date":    []string{"2024-12-01"},
				"launch_pad_id":  []string{"lp-123"},
				"destination_id": []string{"dest-456"},
			},
			expected: &bookingsv1.ListBookingsRequest{
				Filters: bookingsv1.ListBookingsFilters{
					LaunchDate:    stringPtr("2024-12-01"),
					LaunchPadID:   stringPtr("lp-123"),
					DestinationID: stringPtr("dest-456"),
				},
				Pagination: bookingsv1.Pagination{
					Offset: 5,
					Limit:  20,
				},
			},
			expectedErr: nil,
		},
		{
			name: "default parameters",
			params: url.Values{
				"launch_date":    []string{"2024-12-01"},
				"launch_pad_id":  []string{"lp-123"},
				"destination_id": []string{"dest-456"},
			},
			expected: &bookingsv1.ListBookingsRequest{
				Filters: bookingsv1.ListBookingsFilters{
					LaunchDate:    stringPtr("2024-12-01"),
					LaunchPadID:   stringPtr("lp-123"),
					DestinationID: stringPtr("dest-456"),
				},
				Pagination: bookingsv1.Pagination{
					Offset: 0,
					Limit:  10, // Default limit
				},
			},
			expectedErr: nil,
		},
		{
			name: "negative offset",
			params: url.Values{
				"offset": []string{"-1"},
			},
			expected:    nil,
			expectedErr: errors.New("invalid offset"),
		},
		{
			name: "invalid offset",
			params: url.Values{
				"offset": []string{"abc"},
			},
			expected:    nil,
			expectedErr: errors.New("unable to parse offset"),
		},
		{
			name: "negative limit",
			params: url.Values{
				"limit": []string{"-5"},
			},
			expected:    nil,
			expectedErr: errors.New("invalid limit"),
		},
		{
			name: "invalid limit",
			params: url.Values{
				"limit": []string{"xyz"},
			},
			expected:    nil,
			expectedErr: errors.New("unable to parse limit"),
		},
		{
			name: "missing launch_date",
			params: url.Values{
				"offset": []string{"0"},
			},
			expected: &bookingsv1.ListBookingsRequest{
				Pagination: bookingsv1.Pagination{
					Offset: 0,
					Limit:  10, // Default limit
				},
			},
			expectedErr: nil,
		},
		{
			name: "multiple launch_dates",
			params: url.Values{
				"launch_date": []string{"2024-12-01", "2025-01-01"},
				"offset":      []string{"0"},
			},
			expected: &bookingsv1.ListBookingsRequest{
				Filters: bookingsv1.ListBookingsFilters{
					LaunchDate: stringPtr("2024-12-01"), // It takes the first value
				},
				Pagination: bookingsv1.Pagination{
					Offset: 0,
					Limit:  10, // Default limit
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := createListBookingsFromQueryParams(tt.params)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func parseDate(dateStr string) time.Time {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		panic(err)
	}
	return parsedDate
}
