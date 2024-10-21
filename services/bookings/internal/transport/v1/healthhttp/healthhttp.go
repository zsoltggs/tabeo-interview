package healthhttp

import (
	"encoding/json"
	"net/http"

	bookingsv1 "github.com/zsoltggs/tabeo-interview/services/bookings/pkg/bookings/v1"

	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -package=mocks -destination=../../../mocks/healthhttp.go github.com/zsoltggs/tabeo-interview/services/bookings/internal/transport/v1/healthhttp HealthCheckable
type HealthCheckable interface {
	Health() error
}

type HealthHTTP interface {
	HttpHandler(response http.ResponseWriter, request *http.Request)
}

type healthService struct {
	svc HealthCheckable
}

func New(svc HealthCheckable) HealthHTTP {
	return &healthService{
		svc: svc,
	}
}

func (h healthService) HttpHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	err := h.svc.Health()
	if err != nil {
		log.WithError(err).Infof("service is unhealthy")
		response.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	resp := bookingsv1.HealthResponse{
		Status: "OK",
	}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		log.WithError(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	_, err = response.Write(respJSON)
	if err != nil {
		log.WithError(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}
