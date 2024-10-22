package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/transport/v1/bookingshttp"

	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/transport"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/transport/v1/healthhttp"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type httpTransport struct {
	httpServer  *http.Server
	healthSvc   healthhttp.HealthHTTP
	bookingsSvc bookingshttp.BookingsHTTP
}

func NewHTTP(healthSvc healthhttp.HealthHTTP, bookingsSvc bookingshttp.BookingsHTTP) transport.Transport {
	return &httpTransport{
		healthSvc:   healthSvc,
		bookingsSvc: bookingsSvc,
		httpServer:  &http.Server{},
	}
}

func (h *httpTransport) Serve(httpPort int) error {
	port := fmt.Sprintf(":%d", httpPort)
	log.Infof("about to start server on port %s", port)
	router := mux.NewRouter()
	router.HandleFunc("/health", h.healthSvc.HttpHandler).
		Methods("GET")
	router.HandleFunc("/bookings", h.bookingsSvc.ListBookings).
		Methods("GET")
	router.HandleFunc("/bookings", h.bookingsSvc.CreateBooking).
		Methods("POST")
	router.HandleFunc("/bookings/{booking-id}", h.bookingsSvc.DeleteBooking).
		Methods("DELETE")
	h.httpServer.Addr = port
	h.httpServer.Handler = router
	go func() {
		if err := h.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return nil
}

func (h *httpTransport) GracefulStop(ctx context.Context) {
	err := h.httpServer.Shutdown(ctx)
	if err != nil {
		log.WithError(err).Error("unable to shutdown http server")
	}
}
