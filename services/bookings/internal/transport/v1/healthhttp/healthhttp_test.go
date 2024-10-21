package healthhttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	bookingsv1 "github.com/zsoltggs/tabeo-interview/services/bookings/pkg/bookings/v1"

	"github.com/stretchr/testify/assert"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/mocks"
	"go.uber.org/mock/gomock"
)

func Test_Health(t *testing.T) {
	tests := map[string]struct {
		req *http.Request

		mockFunc func(mockHealth *mocks.MockHealthCheckable)

		expectedStatusCode int
		expectedResponse   *bookingsv1.HealthResponse
	}{
		"Invalid http method": {
			req:                httptest.NewRequest(http.MethodPost, "/health", nil),
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		"Unavailable": {
			req: httptest.NewRequest(http.MethodGet, "/health", nil),
			mockFunc: func(mockHealth *mocks.MockHealthCheckable) {
				mockHealth.EXPECT().Health().Return(errors.New("boom"))
			},
			expectedStatusCode: http.StatusServiceUnavailable,
		},
		"Success - Healthy": {
			req: httptest.NewRequest(http.MethodGet, "/health", nil),
			mockFunc: func(mockHealth *mocks.MockHealthCheckable) {
				mockHealth.EXPECT().Health().Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: &bookingsv1.HealthResponse{
				Status: "OK",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockHealth := mocks.NewMockHealthCheckable(ctrl)

			if tc.mockFunc != nil {
				tc.mockFunc(mockHealth)
			}
			handler := New(mockHealth)

			rec := httptest.NewRecorder()
			handler.HttpHandler(rec, tc.req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			if tc.expectedStatusCode == http.StatusOK {
				expectedJSON, _ := json.Marshal(tc.expectedResponse)
				assert.JSONEq(t, string(expectedJSON), rec.Body.String(), "Response body should match")
			}

		})
	}
}
