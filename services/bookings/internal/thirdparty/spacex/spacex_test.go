package spacex

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/models"
)

func TestGetLaunchPadForID(t *testing.T) {
	mockLaunchpads := []Launchpad{
		{Id: "1", Name: "Launchpad 1"},
		{Id: "2", Name: "Launchpad 2"},
	}
	tests := []struct {
		name             string
		launchPadID      string
		mockResponseCode int
		mockResponseBody string
		expectedResult   *Launchpad
		expectedError    error
	}{
		{
			name:             "successful fetch",
			launchPadID:      "1",
			mockResponseCode: http.StatusOK,
			mockResponseBody: `[{"id":"1","name":"Launchpad 1"},{"id":"2","name":"Launchpad 2"}]`,
			expectedResult:   &mockLaunchpads[0],
			expectedError:    nil,
		},
		{
			name:             "launchpad not found",
			launchPadID:      "3",
			mockResponseCode: http.StatusOK,
			mockResponseBody: `[{"id":"1","name":"Launchpad 1"},{"id":"2","name":"Launchpad 2"}]`,
			expectedResult:   nil,
			expectedError:    models.ErrNotFoundLaunchpad,
		},
		{
			name:             "error fetching launchpads",
			launchPadID:      "1",
			mockResponseCode: http.StatusInternalServerError,
			mockResponseBody: ``,
			expectedResult:   nil,
			expectedError:    errors.New("failed to fetch launchpads: status code 500"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock server to simulate HTTP responses
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockResponseCode)
				_, _ = w.Write([]byte(tt.mockResponseBody))
			}))
			defer ts.Close()

			client := ts.Client()
			svc := New(ts.URL, client)

			result, err := svc.GetLaunchPadForID(context.Background(), tt.launchPadID)
			if tt.expectedResult != nil {
				assert.Equal(t, tt.expectedResult.Id, result.Id)
				assert.Nil(t, err)
			} else if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			}
		})
	}
}

func TestGetLaunchesForDate(t *testing.T) {
	mockDate := time.Date(2022, 07, 07, 13, 11, 0, 0, time.UTC)
	mockLaunches := []Launch{
		{
			Name:      "Launch 1",
			DateUTC:   mockDate,
			Launchpad: "launchpad123",
			Success:   true,
		},
		{
			Name:      "Launch 2",
			DateUTC:   mockDate,
			Launchpad: "launchpad456",
			Success:   false,
		},
	}

	tests := []struct {
		name             string
		launchPadID      string
		date             time.Time
		mockResponseCode int
		mockResponseBody string
		expectedResult   []Launch
		expectedError    error
	}{
		{
			name:             "successful fetch",
			launchPadID:      "1",
			date:             mockDate,
			mockResponseCode: http.StatusOK,
			mockResponseBody: `{"docs":[
	{"name":"Launch 1","date_utc":"2022-07-07T13:11:00.000Z","launchpad":"launchpad123","success":true},
	{"name":"Launch 2","date_utc":"2022-07-07T13:11:00.000Z","launchpad":"launchpad456","success":false}
]}`,
			expectedResult: mockLaunches,
			expectedError:  nil,
		},
		{
			name:             "error fetching launches",
			launchPadID:      "1",
			date:             mockDate,
			mockResponseCode: http.StatusInternalServerError,
			mockResponseBody: ``,
			expectedResult:   nil,
			expectedError:    errors.New("failed to fetch launches: status code 500"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock server to simulate HTTP responses
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockResponseCode)
				_, _ = w.Write([]byte(tt.mockResponseBody))
			}))
			defer ts.Close()

			client := ts.Client()
			svc := New(ts.URL, client)

			result, err := svc.GetLaunchesForDate(context.Background(), tt.launchPadID, tt.date)
			if tt.expectedResult != nil {
				assert.Equal(t, tt.expectedResult, result)
				assert.Nil(t, err)
			} else if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			}
		})
	}
}
