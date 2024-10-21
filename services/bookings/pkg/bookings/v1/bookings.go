package v1

type CreateBookingRequest struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Gender        string `json:"gender"`
	BirthDay      string `json:"birthday"`
	LaunchPadID   string `json:"launch_pad_id"`
	DestinationID string `json:"destination_id"`
	LaunchDate    string `json:"launch_date"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

type ListBookingsRequest struct {
	Filters    ListBookingsFilters `json:"filters"`
	Pagination Pagination          `json:"pagination"`
}

type ListBookingsFilters struct {
	LaunchDate    *string `json:"launch_date"`
	LaunchPadID   *string `json:"launch_pad_id"`
	DestinationID *string `json:"destination_id"`
}

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
