package spacex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/thirdparty/spacex/smodels"

	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/models"
)

//go:generate mockgen -package=mocks -destination=../../mocks/spacex.go github.com/zsoltggs/tabeo-interview/services/bookings/internal/thirdparty/spacex SpaceXService
type SpaceXService interface {
	GetLaunchPadForID(ctx context.Context, launchPadID string) (*smodels.Launchpad, error)
	GetLaunchesForDate(ctx context.Context, launchPadID string, date time.Time) ([]smodels.Launch, error)
}

type service struct {
	baseURL string
	client  *http.Client
}

func New(baseURL string, client *http.Client) SpaceXService {
	return &service{
		baseURL: baseURL,
		client:  client,
	}
}

func (s *service) GetLaunchPadForID(ctx context.Context, launchPadID string) (*smodels.Launchpad, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/launchpads", s.baseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create get request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to make request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch launchpads: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var launchpads []smodels.Launchpad
	err = json.Unmarshal(body, &launchpads)
	if err != nil {
		return nil, err
	}

	for _, pad := range launchpads {
		if pad.ID == launchPadID {
			return &pad, nil
		}
	}

	return nil, models.ErrNotFoundLaunchpad
}

func (s *service) GetLaunchesForDate(ctx context.Context, launchPadID string, date time.Time) ([]smodels.Launch, error) {
	startDate := fmt.Sprintf("%sT00:00:00.000Z", date.Format("2006-01-02")) // Start of the day
	endDate := fmt.Sprintf("%sT23:59:59.999Z", date.Format("2006-01-02"))   // End of the day

	queryRequest := smodels.LaunchQueryRequest{
		Query: smodels.LaunchQuery{
			Launchpad: launchPadID,
			DateUTC: smodels.DateRange{
				Gte: startDate,
				Lt:  endDate,
			},
		},
		Options: smodels.LaunchQueryOptions{
			Limit: 5,
		},
	}

	reqBody, err := json.Marshal(queryRequest)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal request: %w", err)
	}
	url := fmt.Sprintf("%s/launches/query", s.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("unable to create get request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to execute request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch launches: status code %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var launchesResponse struct {
		Docs []smodels.Launch `json:"docs"`
	}
	err = json.Unmarshal(respBody, &launchesResponse)
	if err != nil {
		return nil, err
	}

	return launchesResponse.Docs, nil
}
