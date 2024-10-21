package availability

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/mocks"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/thirdparty/spacex"
	"go.uber.org/mock/gomock"
)

func TestIsDateAvailable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSpaceXService := mocks.NewMockSpaceXService(ctrl)
	svc := New(mockSpaceXService)

	const launchPadID = "5e9e4501f509094ba4566f84"
	launch := spacex.Launch{
		Name:      "Starlink 4-21 (v1.5)",
		DateUTC:   time.Date(2022, 07, 07, 13, 11, 00, 0, time.UTC),
		Launchpad: launchPadID,
		Success:   true,
	}
	tests := []struct {
		name          string
		launchPadID   string
		date          time.Time
		mockSetup     func()
		expected      bool
		expectedError error
	}{
		{
			name:        "Available: Valid launch pad with no launches",
			launchPadID: launchPadID,
			date:        time.Date(2022, 7, 7, 0, 0, 0, 0, time.UTC),
			mockSetup: func() {
				mockSpaceXService.EXPECT().
					GetLaunchPadForID(gomock.Any(), launchPadID).
					Return(&spacex.Launchpad{}, nil)
				mockSpaceXService.EXPECT().
					GetLaunchesForDate(gomock.Any(), launchPadID, time.Date(2022, 7, 7, 0, 0, 0, 0, time.UTC)).
					Return(nil, nil)
			},
			expected:      true,
			expectedError: nil,
		},
		{
			name:        "Not Available: Valid launch pad with launches",
			launchPadID: launchPadID,
			date:        time.Date(2022, 7, 7, 0, 0, 0, 0, time.UTC),
			mockSetup: func() {
				mockSpaceXService.EXPECT().
					GetLaunchPadForID(gomock.Any(), launchPadID).
					Return(&spacex.Launchpad{}, nil)
				mockSpaceXService.EXPECT().
					GetLaunchesForDate(gomock.Any(), launchPadID, time.Date(2022, 7, 7, 0, 0, 0, 0, time.UTC)).
					Return([]spacex.Launch{launch}, nil)
			},
			expected:      false,
			expectedError: nil,
		},
		{
			name:        "Invalid launch pad ID",
			launchPadID: "invalid_launchpad_id",
			date:        time.Date(2022, 7, 7, 0, 0, 0, 0, time.UTC),
			mockSetup: func() {
				mockSpaceXService.EXPECT().
					GetLaunchPadForID(gomock.Any(), "invalid_launchpad_id").
					Return(nil, errors.New("not found"))
			},
			expected:      false,
			expectedError: errors.New("unable to get launch pad for ID: not found"),
		},
		{
			name:        "Error getting launches",
			launchPadID: launchPadID,
			date:        time.Date(2022, 7, 7, 0, 0, 0, 0, time.UTC),
			mockSetup: func() {
				mockSpaceXService.EXPECT().
					GetLaunchPadForID(gomock.Any(), launchPadID).
					Return(&spacex.Launchpad{}, nil)
				mockSpaceXService.EXPECT().
					GetLaunchesForDate(gomock.Any(), launchPadID, time.Date(2022, 7, 7, 0, 0, 0, 0, time.UTC)).
					Return(nil, errors.New("internal server error"))
			},
			expected:      false,
			expectedError: errors.New("unable to get launches: internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			available, err := svc.IsDateAvailable(context.Background(), tt.launchPadID, tt.date)

			assert.Equal(t, tt.expected, available)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
