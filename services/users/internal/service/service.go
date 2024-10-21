package service

import (
	"context"
	"fmt"

	"github.com/jonboulle/clockwork"
	"github.com/zsoltggs/tabeo-interview/services/users/internal/database"
	"github.com/zsoltggs/tabeo-interview/services/users/internal/models"
)

//go:generate mockgen -package=mocks -destination=../mocks/service.go github.com/zsoltggs/tabeo-interview/services/users/internal/service Service
type Service interface {
	CreateUser(ctx context.Context, user models.CreateUser) (*models.User, error)
}

type service struct {
	db    database.Database
	clock clockwork.Clock
}

func New(db database.Database,
	clock clockwork.Clock) Service {
	return &service{
		db:    db,
		clock: clock,
	}
}

func (s *service) CreateUser(ctx context.Context, user models.CreateUser) (*models.User, error) {
	now := s.clock.Now()
	// Call db
	if err != nil {
		return nil, fmt.Errorf("unable to hash password: %w", err)
	}
	err = s.db.Create(ctx, models.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Nickname:  user.NickName,
		Password:  password,
		Email:     user.Email,
		Country:   user.Country,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create user: %w", err)
	}
	// Call db to get
	createdUser, err := s.db.GetByID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to get user by id: %w", err)
	}
	return createdUser, nil
}
