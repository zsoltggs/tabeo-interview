package main

import (
	"context"
	http "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/service"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/service/availability"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/thirdparty/spacex"

	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/transport/v1/bookingshttp"

	v1 "github.com/zsoltggs/tabeo-interview/services/bookings/internal/transport/v1"
	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/transport/v1/healthhttp"

	"github.com/zsoltggs/tabeo-interview/services/bookings/internal/database"

	cli "github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"
)

func main() {
	app := cli.App("bookings-service", "bookings service for launchpads")
	pgConnStr := app.String(cli.StringOpt{
		Name:   "pg-connection-string",
		Desc:   "connection string",
		EnvVar: "PG_CONNECTION_STRING",
		Value:  "postgres://user:password@localhost:5432/mydb?sslmode=disable",
	})

	restPort := app.Int(cli.IntOpt{
		Name:   "rest-port",
		Desc:   "rest api port for health check",
		Value:  8080,
		EnvVar: "REST_PORT",
	})
	spaceXBaseURL := app.String(cli.StringOpt{
		Name:   "spacex-base-url",
		Desc:   "url for spacex base",
		Value:  "https://api.spacexdata.com/v4",
		EnvVar: "SPACEX_BASE_URL",
	})

	app.Action = func() {
		log.Info("starting server")

		ctx, cancel := context.WithCancel(context.Background())
		db, err := database.NewPostgres(ctx, *pgConnStr)
		if err != nil {
			log.WithError(err).Panic("unable to connect to postgres")
		}
		defer db.Close(ctx)

		healthSvc := healthhttp.New(db)
		spacexSvc := spacex.NewCache(
			spacex.New(*spaceXBaseURL, &http.Client{
				Timeout: 10 * time.Second,
			}),
			clockwork.NewRealClock(),
		)
		availabilitySvc := availability.New(spacexSvc)
		svc := service.New(db, availabilitySvc, clockwork.NewRealClock(), uuid.New)
		bookingsSvc := bookingshttp.New(svc)

		httpServer := v1.NewHTTP(healthSvc, bookingsSvc)
		err = httpServer.Serve(*restPort)
		if err != nil {
			log.WithError(err).Panic("unable to start http server")
		}
		defer httpServer.GracefulStop(ctx)

		waitForShutdown(cancel)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.WithError(err).Panic("service stopped")
	}
}

// Graceful shutdown
func waitForShutdown(cancel func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	log.Warn("shutting down")
	cancel()
}
