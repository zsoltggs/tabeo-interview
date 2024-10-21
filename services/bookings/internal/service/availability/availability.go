package availability

import (
	"context"
	"time"
)

type Availability interface {
	// IsDateAvailable checks if date is available for the given launchPadID and Date
	IsDateAvailable(ctx context.Context, launchPadID string, date time.Time) (bool, error)
}
