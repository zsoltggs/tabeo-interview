package spacex

import (
	"context"
	"fmt"
	"time"

	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/thirdparty/spacex/smodels"

	"github.com/jonboulle/clockwork"
)

type cache struct {
	svc   SpaceXService
	clock clockwork.Clock

	padCache    map[string]*smodels.Launchpad
	launchCache map[string][]smodels.Launch
}

func NewCache(svc SpaceXService, clock clockwork.Clock) SpaceXService {
	return &cache{
		svc:         svc,
		clock:       clock,
		padCache:    make(map[string]*smodels.Launchpad),
		launchCache: make(map[string][]smodels.Launch),
	}
}

func (c cache) GetLaunchPadForID(ctx context.Context, launchPadID string) (*smodels.Launchpad, error) {
	val, ok := c.padCache[launchPadID]
	if ok {
		return val, nil
	}
	res, err := c.svc.GetLaunchPadForID(ctx, launchPadID)
	if err != nil {
		return nil, fmt.Errorf("unable to get launch pad: %w", err)
	}
	c.padCache[launchPadID] = res
	return res, nil
}

func (c cache) GetLaunchesForDate(ctx context.Context, launchPadID string, date time.Time) ([]smodels.Launch, error) {
	// Assumption future launches might change therefore we only cache launches in the past
	// We could also apply an expiration here but it is fine for now IMO
	now := c.clock.Now()
	if date.Before(now) {
		// Cache
		res, ok := c.launchCache[toLaunchKey(launchPadID, date)]
		if ok {
			return res, nil
		}
		launches, err := c.svc.GetLaunchesForDate(ctx, launchPadID, date)
		if err != nil {
			return nil, fmt.Errorf("unable to get launches: %w", err)
		}
		c.launchCache[toLaunchKey(launchPadID, date)] = launches
		return launches, nil
	} else {
		return c.svc.GetLaunchesForDate(ctx, launchPadID, date)
	}
}

func toLaunchKey(launchPadID string, date time.Time) string {
	return fmt.Sprintf("%s_%s", launchPadID, date.Format("2006-01-02"))
}
