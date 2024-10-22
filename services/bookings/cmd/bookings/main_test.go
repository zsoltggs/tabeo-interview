package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/thirdparty/spacex"
)

func isE2ETestEnabled(t *testing.T) {
	boolValue, err := strconv.ParseBool(os.Getenv("E2E_TEST_ENABLED"))
	if err != nil {
		t.Skip("E2E_TEST_ENABLED environment variable not set")
	}
	if !boolValue {
		t.Skip("E2E_TEST_ENABLED is not true")
	}
}

func Test_SpaceX_GetLaunchPadForID(t *testing.T) {
	isE2ETestEnabled(t)

	svc := spacex.New("https://api.spacexdata.com/v4", &http.Client{
		Timeout: 10 * time.Second,
	})
	ctx := context.Background()
	validLunchPadID := "5e9e4501f509094ba4566f84"
	res, err := svc.GetLaunchPadForID(ctx, validLunchPadID)
	require.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.Name, "CCSFS SLC 40")
}

func Test_SpaceX_GetLaunchesForDate(t *testing.T) {
	isE2ETestEnabled(t)

	validLunchPadID := "5e9e4501f509094ba4566f84"
	expected := spacex.Launch{
		Name:      "Starlink 4-21 (v1.5)",
		DateUTC:   time.Date(2022, 07, 07, 13, 11, 00, 0, time.UTC),
		Launchpad: validLunchPadID,
		Success:   true,
	}

	svc := spacex.New("https://api.spacexdata.com/v4", &http.Client{
		Timeout: 10 * time.Second,
	})
	ctx := context.Background()

	date := time.Date(2022, 07, 07, 1, 1, 1, 1, time.UTC)
	res, err := svc.GetLaunchesForDate(ctx, validLunchPadID, date)
	require.NoError(t, err)
	assert.NotNil(t, res)
	require.Len(t, res, 1)
	assert.Equal(t, expected, res[0])
}

//TODO Add E2E tests for the svc (HTTP)
