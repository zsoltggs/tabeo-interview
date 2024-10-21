package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

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

	app.Action = func() {
		log.Info("starting server")

		ctx, cancel := context.WithCancel(context.Background())
		db, err := database.NewPostgres(ctx, *pgConnStr)
		if err != nil {
			log.WithError(err).Panic("unable to connect to postgres")
		}
		defer db.Close(ctx)

		healthSvc := healthhttp.New(db)
		bookingsSvc := bookingshttp.New()

		http := v1.NewHTTP(healthSvc, bookingsSvc)
		err = http.Serve(*restPort)
		if err != nil {
			log.WithError(err).Panic("unable to start http server")
		}
		defer http.GracefulStop(ctx)

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
