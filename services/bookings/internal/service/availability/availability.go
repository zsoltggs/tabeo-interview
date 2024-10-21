package availability

import (
	"context"
	"fmt"
	"time"

	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/thirdparty/spacex"
)

//go:generate mockgen -package=mocks -destination=../../mocks/availability.go github.com/zsoltggs/tabeo-interview/services/bookings/internal/service/availability  Availability
type Availability interface {
	// IsDateAvailable checks if date is available for the given launchPadID and Date
	IsDateAvailable(ctx context.Context, launchPadID string, date time.Time) (bool, error)
}

type service struct {
	spacexSvc spacex.SpaceXService
}

func New(spacexSvc spacex.SpaceXService) Availability {
	return &service{
		spacexSvc: spacexSvc,
	}
}

func (s service) IsDateAvailable(ctx context.Context, launchPadID string, date time.Time) (bool, error) {
	// Validate that the launch pad is valid
	_, err := s.spacexSvc.GetLaunchPadForID(ctx, launchPadID)
	if err != nil {
		return false, fmt.Errorf("unable to get launch pad for ID: %w", err)
	}
	// Get all launches for the date
	launches, err := s.spacexSvc.GetLaunchesForDate(ctx, launchPadID, date)
	if err != nil {
		return false, fmt.Errorf("unable to get launches: %w", err)
	}
	if len(launches) != 0 {
		return false, nil
	}
	// If no launches are for the date it means that the date is available
	return true, nil
}
