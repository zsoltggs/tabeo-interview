package database

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/models"
)

func isIntegrationTestEnabled(t *testing.T) {
	boolValue, err := strconv.ParseBool(os.Getenv("INTEGRATION_TEST_ENABLED"))
	if err != nil {
		t.Skip("INTEGRATION_TEST_ENABLED environment variable not set")
	}
	if !boolValue {
		t.Skip("INTEGRATION_TEST_ENABLED is not true")
	}
}

func setupTestDB(t *testing.T) Database {
	isIntegrationTestEnabled(t)
	connectionStr := "postgres://myuser:mypassword@localhost:5432/bookings?sslmode=disable"
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connectionStr)
	assert.NoError(t, err)
	defer pool.Close()

	_, err = pool.Exec(context.Background(), "TRUNCATE TABLE bookings")
	assert.NoError(t, err)

	db, err := NewPostgres(ctx, connectionStr)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	return db // Use the concrete type for testing purposes
}

func TestCreateBooking(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close(context.Background())

	id := uuid.New()
	now := time.Now()

	booking := models.Booking{
		ID:            id,
		FirstName:     "John",
		LastName:      "Doe",
		Gender:        "Male",
		Birthday:      now.AddDate(-25, 0, 0),
		LaunchPadID:   "LP-001",
		DestinationID: "DS-001",
		LaunchDate:    now.AddDate(0, 1, 0),
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	err := db.Create(context.Background(), booking)
	assert.NoError(t, err)

	// Retrieve the booking and verify its details
	savedBooking, err := db.GetByID(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, booking.ID, savedBooking.ID)
	assert.Equal(t, booking.FirstName, savedBooking.FirstName)
}

func TestDeleteBooking(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close(context.Background())

	id := uuid.New()
	now := time.Now()

	booking := models.Booking{
		ID:            id,
		FirstName:     "Jane",
		LastName:      "Doe",
		Gender:        "Female",
		Birthday:      now.AddDate(-30, 0, 0),
		LaunchPadID:   "LP-002",
		DestinationID: "DS-002",
		LaunchDate:    now.AddDate(0, 2, 0),
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	err := db.Create(context.Background(), booking)
	assert.NoError(t, err)

	err = db.Delete(context.Background(), id)
	assert.NoError(t, err)

	_, err = db.GetByID(context.Background(), id)
	assert.Error(t, err)
}

func TestListBookings(t *testing.T) {
	// Set up the test DB
	db := setupTestDB(t)
	defer db.Close(context.Background())

	// Seed the database with test bookings
	now := time.Now()

	// Create multiple bookings for testing
	for i := 0; i < 5; i++ {
		booking := models.Booking{
			ID:            uuid.New(),
			FirstName:     fmt.Sprintf("TestFirstName-%d", i),
			LastName:      fmt.Sprintf("TestLastName-%d", i),
			Gender:        "Other",
			Birthday:      now.AddDate(-20, 0, 0),
			LaunchPadID:   fmt.Sprintf("LP-00%d", i+1),
			DestinationID: fmt.Sprintf("DS-00%d", i+1),
			LaunchDate:    now.AddDate(0, i, 0),
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		err := db.Create(context.Background(), booking)
		assert.NoError(t, err, "Failed to create booking")
	}

	// Define pagination and filters
	pagination := models.Pagination{
		Offset: 0,
		Limit:  3, // Fetch 3 bookings at a time
	}

	filters := models.Filters{
		LaunchDate:    nil,
		LaunchPadID:   nil,
		DestinationID: nil,
	}

	// Fetch bookings with no filters
	bookings, err := db.List(context.Background(), pagination, filters)
	assert.NoError(t, err, "Failed to list bookings without filters")
	assert.Len(t, bookings, 3, "Expected 3 bookings in the first batch")

	// Test filter by LaunchPadID
	launchPadID := "LP-001"
	filters.LaunchPadID = &launchPadID

	bookings, err = db.List(context.Background(), pagination, filters)
	assert.NoError(t, err, "Failed to list bookings by LaunchPadID")
	assert.Len(t, bookings, 1, "Expected 1 booking with the specified LaunchPadID")
	assert.Equal(t, "LP-001", bookings[0].LaunchPadID, "Incorrect LaunchPadID returned")

	// Test filter by DestinationID
	destinationID := "DS-002"
	filters.LaunchPadID = nil // Reset LaunchPadID filter
	filters.DestinationID = &destinationID

	bookings, err = db.List(context.Background(), pagination, filters)
	assert.NoError(t, err, "Failed to list bookings by DestinationID")
	assert.Len(t, bookings, 1, "Expected 1 booking with the specified DestinationID")
	assert.Equal(t, "DS-002", bookings[0].DestinationID, "Incorrect DestinationID returned")

	// Test pagination by fetching the next set of bookings
	pagination.Offset = 3
	bookings, err = db.List(context.Background(), models.Pagination{
		Offset: 3, // Skip the first 3 bookings
		Limit:  3,
	}, models.Filters{})
	assert.NoError(t, err, "Failed to list bookings with pagination")
	assert.Len(t, bookings, 2, "Expected 2 bookings in the second batch")
}

func TestHealth(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close(context.Background())

	err := db.Health()
	assert.NoError(t, err, "Database health check failed")
}
