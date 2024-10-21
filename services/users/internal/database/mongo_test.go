package database

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/zsoltggs/tabeo-interview/services/users/internal/models"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoConnectionString = "mongodb://localhost:27017"

func truncateTables(t *testing.T, mongoURL, databaseName string) {
	var err error

	opts := options.Client().ApplyURI(mongoURL)
	client, err := mongo.Connect(context.Background(), opts)
	assert.Nil(t, err)
	err = client.Database(databaseName).Drop(context.Background())
	assert.Nil(t, err)
}

func isIntegrationTestEnabled(t *testing.T) {
	boolValue, err := strconv.ParseBool(os.Getenv("INTEGRATION_TEST_ENABLED"))
	if err != nil {
		t.Skip("INTEGRATION_TEST_ENABLED environment variable not set")
	}
	if !boolValue {
		t.Skip("INTEGRATION_TEST_ENABLED is not true")
	}
}

func initDB(t *testing.T, connectionString string, databaseName string) Database {
	isIntegrationTestEnabled(t)
	dbName := t.Name()
	truncateTables(t, connectionString, databaseName)
	ctx := context.Background()
	mgo, err := NewMongo(ctx, connectionString, dbName)
	require.NoError(t, err)
	return mgo
}

func Test_MongoDB_Create(t *testing.T) {
	db := initDB(t, mongoConnectionString, t.Name())
	ctx := context.Background()
	defer db.Close(ctx)

	err := db.Create(ctx, createTestUser())
	require.NoError(t, err)
}

func Test_MongoDB_GetByID(t *testing.T) {
	db := initDB(t, mongoConnectionString, t.Name())
	ctx := context.Background()
	defer db.Close(ctx)

	usr := createTestUser()
	err := db.Create(ctx, usr)
	require.NoError(t, err)

	result, err := db.GetByID(ctx, usr.ID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, usr, *result)
}

func Test_MongoDB_GetByID_ReturnsNotFound(t *testing.T) {
	db := initDB(t, mongoConnectionString, t.Name())
	ctx := context.Background()
	defer db.Close(ctx)

	result, err := db.GetByID(ctx, uuid.New())
	require.Error(t, err)
	require.Nil(t, result)
	assert.ErrorContains(t, err, ErrNotFound.Error())
}

func Test_MongoDB_Update_ReturnsNotFound(t *testing.T) {
	db := initDB(t, mongoConnectionString, t.Name())
	ctx := context.Background()
	defer db.Close(ctx)

	err := db.Update(ctx, createTestUser())
	require.Error(t, err)
	assert.ErrorContains(t, err, ErrNotFound.Error())
}

func Test_MongoDB_Update(t *testing.T) {
	db := initDB(t, mongoConnectionString, t.Name())
	ctx := context.Background()
	defer db.Close(ctx)

	usr := createTestUser()
	err := db.Create(ctx, usr)
	require.NoError(t, err)

	usr.Email = "asd@gmail.com"
	usr.LastName = "dbf"
	usr.Password = "new-password"
	err = db.Update(ctx, usr)
	require.NoError(t, err)

	result, err := db.GetByID(ctx, usr.ID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, usr, *result)
}

func Test_MongoDB_Delete_ErrNotFound(t *testing.T) {
	db := initDB(t, mongoConnectionString, t.Name())
	ctx := context.Background()
	defer db.Close(ctx)

	err := db.Delete(ctx, uuid.New())
	require.Error(t, err)
	assert.ErrorContains(t, err, ErrNotFound.Error())
}

func Test_MongoDB_Delete(t *testing.T) {
	db := initDB(t, mongoConnectionString, t.Name())
	ctx := context.Background()
	defer db.Close(ctx)

	usr := createTestUser()
	err := db.Create(ctx, usr)
	require.NoError(t, err)

	result, err := db.GetByID(ctx, usr.ID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, usr, *result)

	err = db.Delete(ctx, usr.ID)
	require.NoError(t, err)

	res, err := db.GetByID(ctx, usr.ID)
	require.Error(t, err)
	require.Nil(t, res)
	assert.ErrorContains(t, err, ErrNotFound.Error())
}

func Test_MongoDB_List(t *testing.T) {
	db := initDB(t, mongoConnectionString, t.Name())
	ctx := context.Background()
	defer db.Close(ctx)

	records := addMockDataForList(t, db)
	tests := map[string]struct {
		pagination models.Pagination
		filters    models.Filters
		expected   []models.User
	}{
		"NoFilters": {
			pagination: models.Pagination{
				Offset: 0,
				Limit:  100,
			},
			filters:  models.Filters{},
			expected: records,
		},
		"UK Filter": {
			pagination: models.Pagination{
				Offset: 0,
				Limit:  100,
			},
			filters: models.Filters{
				Country: toPTR("UK"),
			},
			expected: records[0:10],
		},
		"US Filter": {
			pagination: models.Pagination{
				Offset: 0,
				Limit:  100,
			},
			filters: models.Filters{
				Country: toPTR("US"),
			},
			expected: records[10:],
		},
		"Email filter": {
			pagination: models.Pagination{
				Offset: 0,
				Limit:  100,
			},
			filters: models.Filters{
				Email: toPTR("alice@0.com"),
			},
			expected: []models.User{records[0]},
		},
		"Email & Country filter": {
			pagination: models.Pagination{
				Offset: 0,
				Limit:  100,
			},
			filters: models.Filters{
				Country: toPTR("UK"),
				Email:   toPTR("alice@0.com"),
			},
			expected: []models.User{records[0]},
		},
		"Email & Country filter -> not found": {
			pagination: models.Pagination{
				Offset: 0,
				Limit:  100,
			},
			filters: models.Filters{
				Country: toPTR("US"),
				Email:   toPTR("alice@0.com"),
			},
			expected: nil,
		},
		"Test pagination page 0": {
			pagination: models.Pagination{
				Offset: 0,
				Limit:  10,
			},
			expected: records[:10],
		},
		"Test pagination page 1": {
			pagination: models.Pagination{
				Offset: 10,
				Limit:  10,
			},
			expected: records[10:],
		},
		"Test pagination page 0 with filters": {
			pagination: models.Pagination{
				Offset: 0,
				Limit:  5,
			},
			filters: models.Filters{
				Country: toPTR("UK"),
			},
			expected: records[0:5],
		},
		"Test pagination page 1 with filters": {
			pagination: models.Pagination{
				Offset: 5,
				Limit:  5,
			},
			filters: models.Filters{
				Country: toPTR("UK"),
			},
			expected: records[5:10],
		},
		"Test pagination page 2 with filters -> no more pages": {
			pagination: models.Pagination{
				Offset: 10,
				Limit:  5,
			},
			filters: models.Filters{
				Country: toPTR("UK"),
			},
			expected: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			results, err := db.List(ctx, tc.pagination, tc.filters)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, results)
		})
	}
}

func Test_MongoDB_List_US_Filter(t *testing.T) {
	db := initDB(t, mongoConnectionString, t.Name())
	ctx := context.Background()
	defer db.Close(ctx)

	records := addMockDataForList(t, db)

	results, err := db.List(ctx, models.Pagination{
		Offset: 0,
		Limit:  100,
	}, models.Filters{
		Country: toPTR("UK"),
	})
	require.NoError(t, err)
	require.NotNil(t, results)
	assert.Equal(t, records[0:10], results)
}

func addMockDataForList(t *testing.T, db Database) []models.User {
	ctx := context.Background()
	var records []models.User
	now := time.Now().UTC().Truncate(time.Millisecond)
	for i := 0; i < 10; i++ {
		now = now.Add(time.Second * time.Duration(i))
		name := fmt.Sprintf("%d-Alice", i)
		email := fmt.Sprintf("alice@%d.com", i)
		usr := createTestUserWithParams(name, "UK", email, now)
		err := db.Create(ctx, usr)
		require.NoError(t, err)
		records = append(records, usr)
	}
	for i := 0; i < 10; i++ {
		now = now.Add(time.Second * time.Duration(i))
		name := fmt.Sprintf("%d-BobUS", i)
		email := fmt.Sprintf("bob@%d.com", i)
		usr := createTestUserWithParams(name, "US", email, now)
		err := db.Create(ctx, usr)
		require.NoError(t, err)
		records = append(records, usr)
	}
	return records
}

func toPTR(s string) *string {
	return &s
}

func createTestUser() models.User {
	return models.User{
		ID:        uuid.New(),
		FirstName: "Alice",
		LastName:  "Bob",
		Nickname:  "AB",
		Password:  "1237895798kjshf",
		Email:     "alice@bob.com",
		Country:   "UK",
		CreatedAt: time.Now().UTC().Truncate(time.Millisecond),
		UpdatedAt: time.Now().UTC().Truncate(time.Millisecond),
	}
}

func createTestUserWithParams(name, country, email string, ts time.Time) models.User {
	return models.User{
		ID:        uuid.New(),
		FirstName: name,
		LastName:  "Bob",
		Nickname:  "AB",
		Password:  "1237895798kjshf",
		Email:     email,
		Country:   country,
		CreatedAt: ts,
		UpdatedAt: time.Now().UTC().Truncate(time.Millisecond),
	}
}
