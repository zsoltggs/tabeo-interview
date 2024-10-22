package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	bookingsv1 "github.com/zsoltggs/tabeo-interview/services/bookings/pkg/bookings/v1"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/thirdparty/spacex"
)

const (
	serviceBaseURL   = "http://localhost:9999"
	spaceXBaseURL    = "https://api.spacexdata.com/v4"
	validLaunchPadID = "5e9e4501f509094ba4566f84"
)

var launchExistsDate = time.Date(2022, 07, 07, 13, 11, 0, 0, time.UTC)
var launchAvailableDate = time.Date(2022, 07, 8, 1, 1, 1, 1, time.UTC)

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

	svc := spacex.New(spaceXBaseURL, &http.Client{
		Timeout: 10 * time.Second,
	})
	ctx := context.Background()
	res, err := svc.GetLaunchPadForID(ctx, validLaunchPadID)
	require.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.Name, "CCSFS SLC 40")
}

func Test_SpaceX_GetLaunchesForDate(t *testing.T) {
	isE2ETestEnabled(t)

	expected := spacex.Launch{
		Name:      "Starlink 4-21 (v1.5)",
		DateUTC:   launchExistsDate,
		Launchpad: validLaunchPadID,
		Success:   true,
	}

	svc := spacex.New(spaceXBaseURL, &http.Client{
		Timeout: 10 * time.Second,
	})
	ctx := context.Background()

	date := launchExistsDate
	res, err := svc.GetLaunchesForDate(ctx, validLaunchPadID, date)
	require.NoError(t, err)
	assert.NotNil(t, res)
	require.Len(t, res, 1)
	assert.Equal(t, expected, res[0])
}

func Test_Service_E2E(t *testing.T) {
	isE2ETestEnabled(t)

	// Test create booking
	booking := map[string]interface{}{
		"first_name":     "John",
		"last_name":      "Doe",
		"gender":         "male",
		"birthday":       "1990-01-01",
		"launch_pad_id":  validLaunchPadID,
		"destination_id": "some-destination-id",
		"launch_date":    launchAvailableDate.Format("2006-01-02"),
	}

	body, err := json.Marshal(booking)
	assert.NoError(t, err)
	resp, err := http.Post(serviceBaseURL+"/bookings", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode, string(body))
	createBookingResponse := bookingsv1.CreateBookingResponse{}
	err = json.Unmarshal(body, &createBookingResponse)
	require.NoError(t, err)
	assert.Equal(t, "John", createBookingResponse.Booking.FirstName)
	assert.Equal(t, "Doe", createBookingResponse.Booking.LastName)
	assert.Equal(t, "male", createBookingResponse.Booking.Gender)

	// Test List Booking
	listResponse, err := http.Get(serviceBaseURL + "/bookings")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, listResponse.StatusCode)
	defer listResponse.Body.Close()

	// Test delete booking
	deleteURL := serviceBaseURL + "/bookings/" + createBookingResponse.Booking.ID.String()
	req, err := http.NewRequest(http.MethodDelete, deleteURL, nil)
	assert.NoError(t, err)

	client := &http.Client{}
	deleteResp, err := client.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, deleteResp)
	assert.Equal(t, http.StatusNoContent, deleteResp.StatusCode)
}
