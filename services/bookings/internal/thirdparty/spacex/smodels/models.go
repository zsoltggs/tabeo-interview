package smodels

import "time"

// Launch represents the structure of a launch from the SpaceX API
type Launch struct {
	Name      string    `json:"name"`
	DateUTC   time.Time `json:"date_utc"`
	Launchpad string    `json:"launchpad"`
	Success   bool      `json:"success"`
}

// LaunchQueryRequest represents the query for fetching launches
type LaunchQueryRequest struct {
	Query   LaunchQuery        `json:"query"`
	Options LaunchQueryOptions `json:"options"`
}

// LaunchQueryOptions holds the limit option fot the /launches/query endpoint
type LaunchQueryOptions struct {
	Limit int `json:"limit"`
}

// LaunchQuery holds the launchpad and date query for the /launches/query endpoint
type LaunchQuery struct {
	Launchpad string    `json:"launchpad"`
	DateUTC   DateRange `json:"date_utc"`
}

// DateRange defines the date range for querying launches
type DateRange struct {
	Gte string `json:"$gte"` // Greater than or equal to
	Lt  string `json:"$lt"`  // Less than
}

// Launchpad represents the structure of a launchpad in SpaceX API
type Launchpad struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Locality string `json:"locality"`
	Region   string `json:"region"`
	Status   string `json:"status"`
	ID       string `json:"id"`
}
