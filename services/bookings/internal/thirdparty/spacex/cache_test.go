package spacex

import (
	"context"
	"testing"
	"time"

	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/thirdparty/spacex/smodels"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/mocks"
	"go.uber.org/mock/gomock"
)

func Test_Cache_GetLaunchPadForID(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockService := mocks.NewMockSpaceXService(ctrl)
	clock := clockwork.NewFakeClock()
	cachedService := NewCache(mockService, clock)

	launchPadID := "pad-1"
	expectedLaunchPad := &smodels.Launchpad{ID: launchPadID, Name: "Launch Pad 1"}

	// First call should hit the underlying service
	mockService.EXPECT().
		GetLaunchPadForID(context.Background(), launchPadID).
		Return(expectedLaunchPad, nil).Times(1)

	// Call the method
	launchPad, err := cachedService.GetLaunchPadForID(context.Background(), launchPadID)

	assert.NoError(t, err)
	assert.Equal(t, expectedLaunchPad, launchPad)

	// Second call should return cached value
	launchPad, err = cachedService.GetLaunchPadForID(context.Background(), launchPadID)
	assert.NoError(t, err)
	assert.Equal(t, expectedLaunchPad, launchPad)
}

func TestGetLaunchesForDate_CachesResults(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockService := mocks.NewMockSpaceXService(ctrl)
	now := time.Date(2025, 12, 01, 0, 0, 0, 0, time.UTC)
	clock := clockwork.NewFakeClockAt(now)

	cachedService := NewCache(mockService, clock)

	launchPadID := "pad-1"
	date := now.AddDate(0, 0, -10)
	expectedLaunches := []smodels.Launch{
		{Name: "Launch 1"},
		{Name: "Launch 2"},
	}

	// First call should hit the underlying service
	mockService.EXPECT().
		GetLaunchesForDate(context.Background(), launchPadID, date).
		Return(expectedLaunches, nil).Times(1)

	// Call the method
	launches, err := cachedService.GetLaunchesForDate(context.Background(), launchPadID, date)

	assert.NoError(t, err)
	assert.Equal(t, expectedLaunches, launches)

	// Second call should return cached value
	launches, err = cachedService.GetLaunchesForDate(context.Background(), launchPadID, date)
	assert.NoError(t, err)
	assert.Equal(t, expectedLaunches, launches)
}

func TestGetLaunchesForDate_NonCachedFutureDate(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockService := mocks.NewMockSpaceXService(ctrl)
	clock := clockwork.NewFakeClock()
	cachedService := NewCache(mockService, clock)

	launchPadID := "pad-1"
	futureDate := clock.Now().Add(24 * time.Hour)

	// Future date call should hit the underlying service
	mockService.EXPECT().
		GetLaunchesForDate(context.Background(), launchPadID, futureDate).
		Return([]smodels.Launch{}, nil).Times(1)

	// Call the method with a future date
	launches, err := cachedService.GetLaunchesForDate(context.Background(), launchPadID, futureDate)
	assert.NoError(t, err)
	assert.Empty(t, launches)
}

func TestGetLaunchesForDate_ErrorFromService(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockService := mocks.NewMockSpaceXService(ctrl)
	clock := clockwork.NewFakeClock()
	cachedService := NewCache(mockService, clock)

	launchPadID := "pad-1"
	date := time.Date(2024, 12, 01, 0, 0, 0, 0, time.UTC)

	// Simulate an error from the service
	mockService.EXPECT().
		GetLaunchesForDate(context.Background(), launchPadID, date).
		Return(nil, assert.AnError).Times(1)

	// Call the method and expect an error
	launches, err := cachedService.GetLaunchesForDate(context.Background(), launchPadID, date)

	assert.Error(t, err)
	assert.Nil(t, launches)
}
