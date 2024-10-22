package bookingshttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	bookingsv1 "github.com/zsoltggs/tabeo-interview/services/bookings/pkg/bookings/v1"

	"github.com/stretchr/testify/assert"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/mocks"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/models"
	gomock "go.uber.org/mock/gomock"
)

func TestCreateBookingHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	fixedUUID := uuid.MustParse("0aadd991-953d-48d3-a4a8-8e1182a2c723")
	ts := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name           string
		method         string
		body           interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Invalid method",
			method: http.MethodGet,
			body:   nil,
			mockSetup: func() {
			},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "",
		},
		{
			name:   "Invalid JSON body",
			method: http.MethodPost,
			body:   "invalid-body",
			mockSetup: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"bad request"}`,
		},
		{
			name:   "Invalid request required value missing",
			method: http.MethodPost,
			body: bookingsv1.CreateBookingRequest{
				FirstName: "John",
				LastName:  "Doe",
			},
			mockSetup: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "destination id is required"}`,
		},
		{
			name:   "Launch pad not found",
			method: http.MethodPost,
			body: bookingsv1.CreateBookingRequest{
				FirstName:     "John",
				LastName:      "Doe",
				Gender:        "male",
				Birthday:      "1980-01-01",
				LaunchPadID:   "invalid-pad",
				DestinationID: "dest-123",
				LaunchDate:    "2024-12-31",
			},
			mockSetup: func() {
				mockService.EXPECT().
					CreateBooking(gomock.Any(), gomock.Any()).
					Return(nil, models.ErrNotFoundLaunchpad)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error": "launch pad with ID not found"}`,
		},
		{
			name:   "Date unavailable",
			method: http.MethodPost,
			body: bookingsv1.CreateBookingRequest{
				FirstName:     "Jane",
				LastName:      "Doe",
				Gender:        "female",
				Birthday:      "1990-01-01",
				LaunchPadID:   "valid-pad",
				DestinationID: "dest-456",
				LaunchDate:    "2024-12-31",
			},
			mockSetup: func() {
				mockService.EXPECT().
					CreateBooking(gomock.Any(), gomock.Any()).
					Return(nil, models.ErrNotAvailable)
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"error":"date is unavailable"}`,
		},
		{
			name:   "Successful booking",
			method: http.MethodPost,
			body: bookingsv1.CreateBookingRequest{
				FirstName:     "Jane",
				LastName:      "Doe",
				Gender:        "female",
				Birthday:      "1990-01-01",
				LaunchPadID:   "valid-pad",
				DestinationID: "dest-456",
				LaunchDate:    "2024-12-31",
			},
			mockSetup: func() {
				booking := &models.Booking{
					ID:            fixedUUID,
					FirstName:     "Jane",
					LastName:      "Doe",
					Gender:        "female",
					Birthday:      timeDate(1990, 1, 1),
					LaunchPadID:   "valid-pad",
					DestinationID: "dest-456",
					LaunchDate:    timeDate(2024, 12, 31),
					CreatedAt:     ts,
					UpdatedAt:     ts,
				}
				mockService.EXPECT().
					CreateBooking(gomock.Any(), gomock.Any()).
					Return(booking, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: `{"booking":
	{
		"id":"0aadd991-953d-48d3-a4a8-8e1182a2c723",
		"first_name":"Jane",
		"last_name":"Doe",
		"gender":"female",
		"birthday":"1990-01-01",
		"launch_pad_id":"valid-pad",
		"destination_id":"dest-456",
		"launch_date":"2024-12-31",
		"created_at":"2024-01-02T03:04:05Z", 
		"updated_at":"2024-01-02T03:04:05Z"
	}
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			var requestBody []byte
			if tt.body != nil {
				switch v := tt.body.(type) {
				case string:
					requestBody = []byte(v)
				default:
					jsonBody, _ := json.Marshal(v)
					requestBody = jsonBody
				}
			}

			req := httptest.NewRequest(tt.method, "/bookings", bytes.NewReader(requestBody))
			rec := httptest.NewRecorder()

			handler := New(mockService)
			handler.CreateBooking(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			}
		})
	}
}

func TestListBookings(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	fixedUUID := uuid.MustParse("0aadd991-953d-48d3-a4a8-8e1182a2c723")

	handler := bookingsHTTP{
		service: mockService,
	}
	launchDate := timeDate(2024, 01, 02)
	ts := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name           string
		method         string
		queryParams    string
		mockService    func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Method Not Allowed",
			method:         http.MethodPost, // Invalid method
			queryParams:    "",
			mockService:    func() {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "",
		},
		{
			name:           "Bad Request - Invalid Query Params",
			method:         http.MethodGet,
			queryParams:    "?limit=invalid", // Invalid query param
			mockService:    func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "unable to parse limit"}`,
		},
		{
			name:        "Service Error",
			method:      http.MethodGet,
			queryParams: "?offset=2&limit=20",
			mockService: func() {
				mockService.EXPECT().
					ListBookings(gomock.Any(), gomock.Any(), models.Pagination{Offset: 2, Limit: 20}).
					Return(nil, errors.New("some service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "",
		},
		{
			name:        "Successful Response",
			method:      http.MethodGet,
			queryParams: "?launch_date=2024-01-02&launch_pad_id=valid-pad&destination_id=dest-456",
			mockService: func() {
				booking := models.Booking{
					ID:            fixedUUID,
					FirstName:     "Jane",
					LastName:      "Doe",
					Gender:        "female",
					Birthday:      timeDate(1990, 1, 1),
					LaunchPadID:   "valid-pad",
					DestinationID: "dest-456",
					LaunchDate:    timeDate(2024, 12, 31),
					CreatedAt:     ts,
					UpdatedAt:     ts,
				}
				mockService.EXPECT().
					ListBookings(gomock.Any(), models.Filters{
						LaunchDate:    &launchDate,
						LaunchPadID:   toPtr("valid-pad"),
						DestinationID: toPtr("dest-456"),
					}, models.Pagination{Offset: 0, Limit: 10}).
					Return([]models.Booking{
						booking,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{"bookings":[
	{
		"id":"0aadd991-953d-48d3-a4a8-8e1182a2c723",
		"first_name":"Jane",
		"last_name":"Doe",
		"gender":"female",
		"birthday":"1990-01-01",
		"launch_pad_id":"valid-pad",
		"destination_id":"dest-456",
		"launch_date":"2024-12-31",
		"created_at":"2024-01-02T03:04:05Z", 
		"updated_at":"2024-01-02T03:04:05Z"
	}
]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockService()

			req := httptest.NewRequest(tt.method, "/bookings"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			handler.ListBookings(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedBody != "" {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				assert.JSONEq(t, tt.expectedBody, string(body))
			}
		})
	}
}

func timeDate(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func toPtr(s string) *string {
	return &s
}
