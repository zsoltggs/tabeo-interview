package database

import (
	"context"
	"fmt"

	"github.com/zsoltggs/tabeo-interview/services/users/internal/models"

	"github.com/sirupsen/logrus"

	"github.com/globalsign/mgo/bson"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionUsers = "users"
)

type mongoDB struct {
	client   *mongo.Client
	database string
}

func NewMongo(ctx context.Context, connString string, database string) (Database, error) {
	opts := options.Client().ApplyURI(connString).
		SetRetryWrites(true)

	err := opts.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo connection: %w", err)
	}

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo connection: %w", err)
	}

	svc := &mongoDB{client: client, database: database}

	err = svc.ensureIndicies()
	if err != nil {
		return nil, fmt.Errorf("failed to create collection indicies: %w", err)
	}

	return svc, nil
}

func (m *mongoDB) ensureIndicies() error {
	ctx := context.Background()
	session, err := m.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	db := session.Client().Database(m.database)

	indexes := map[string][]mongo.IndexModel{
		collectionUsers: {
			mongo.IndexModel{
				Keys:    bson.M{"id": 1},
				Options: options.Index().SetName("users_idx"),
			},
			mongo.IndexModel{
				Keys:    bson.M{"email": 1},
				Options: options.Index().SetName("users_email_idx"),
			},
			mongo.IndexModel{
				Keys:    bson.M{"country": 1},
				Options: options.Index().SetName("users_country_idx"),
			},
		},
	}

	for collection, indices := range indexes {
		_, err := db.
			Collection(collection).
			Indexes().
			CreateMany(ctx, indices)
		if err != nil {
			return fmt.Errorf("unable to create database indices for collection %q: %w", collection, err)
		}
	}

	return nil
}

func (m *mongoDB) Health() error {
	return m.client.Ping(context.Background(), nil)
}

func (m *mongoDB) Close(ctx context.Context) {
	err := m.client.Disconnect(ctx)
	if err != nil {
		logrus.WithError(err).Error("unable to disconnect from mongo")
	}
}

func (m *mongoDB) Create(ctx context.Context, user models.User) error {
	session, err := m.client.StartSession()
	if err != nil {
		return fmt.Errorf("error creating session: %w", err)
	}

	defer session.EndSession(ctx)
	db := session.Client().Database(m.database)

	_, err = db.Collection(collectionUsers).UpdateOne(ctx,
		bson.M{"id": user.ID},
		bson.M{"$set": user},
		options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

func (m *mongoDB) Update(ctx context.Context, user models.User) error {
	session, err := m.client.StartSession()
	if err != nil {
		return fmt.Errorf("error creating session: %w", err)
	}

	defer session.EndSession(ctx)
	db := session.Client().Database(m.database)

	res, err := db.Collection(collectionUsers).UpdateOne(ctx,
		bson.M{"id": user.ID},
		bson.M{"$set": user},
		options.Update().SetUpsert(false))
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

func (m *mongoDB) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := m.client.Database(m.database).
		Collection(collectionUsers).
		DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("could not delete user by id: %w", err)
	}
	if res.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (m *mongoDB) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	res := m.client.Database(m.database).
		Collection(collectionUsers).
		FindOne(ctx, bson.M{"id": id})
	switch {
	case res.Err() == mongo.ErrNoDocuments:
		return nil, ErrNotFound
	case res.Err() != nil:
		return nil, fmt.Errorf("could not get user: %w", res.Err())
	}

	var u models.User
	err := res.Decode(&u)
	if err != nil {
		return nil, fmt.Errorf("could not decode user: %w", err)
	}
	return &u, nil
}

func (m *mongoDB) List(ctx context.Context, pagination models.Pagination, userFilters models.Filters) ([]models.User, error) {
	session, err := m.client.StartSession()
	if err != nil {
		return nil, err
	}

	defer session.EndSession(ctx)
	db := session.Client().Database(m.database)

	var filters []bson.M
	if fromPtr(userFilters.Email) != "" {
		filters = append(filters, bson.M{"email": *userFilters.Email})
	}
	if fromPtr(userFilters.Country) != "" {
		filters = append(filters, bson.M{"country": *userFilters.Country})
	}

	sorting := bson.M{"createdat": 1}
	query := bson.M{}
	if len(filters) != 0 {
		query = bson.M{"$and": filters}
	}
	find, err := db.Collection(collectionUsers).Find(
		ctx,
		query,
		options.Find().
			SetSkip(int64(pagination.Offset)).
			SetLimit(int64(pagination.Limit)).
			SetSort(sorting),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to list users: %w", err)
	}

	var results []models.User
	err = find.All(ctx, &results)
	if err != nil {
		return nil, fmt.Errorf("unable to find all results: %w", err)
	}

	return results, nil
}

func fromPtr(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}
