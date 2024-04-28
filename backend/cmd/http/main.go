package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/rampa2510/contracts-poc/config"
	"github.com/rampa2510/contracts-poc/db/seed"
	"github.com/rampa2510/contracts-poc/internal/api"
	"github.com/rampa2510/contracts-poc/internal/api/middleware"
	"github.com/rampa2510/contracts-poc/internal/storage"
	"github.com/rampa2510/contracts-poc/pkg/shutdown"
)

func main() {
	// setup exit code for graceful shutdown
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	// load config
	env, err := config.LoadConfig()
	if err != nil {
		slog.Error("error: %v", err)
		exitCode = 1
		return
	}

	cleanup, err := run(env)

	defer cleanup()
	if err != nil {
		slog.Error("error: %v", err)
		exitCode = 1
		return
	}

	shutdown.Gracefully()
}

func run(env config.EnvVars) (func(), error) {
	app, cleanup, err := buildServer(env)
	if err != nil {
		return nil, err
	}

	// start the server
	go func() {
		if err := app.ListenAndServe(); err != nil {
			slog.Error("Error while starting server - %v", err)
		}
	}()

	return func() {
		cleanup()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := app.Shutdown(ctx); err != nil {
			slog.Error("Error while shutting down server - %v", err)
		}
	}, nil
}

func buildServer(env config.EnvVars) (*http.Server, func(), error) {
	slog.Info("Initializing DB connection")
	// init the storage
	db, err := storage.BootstrapSqlite3(env.DB_FILE, 10*time.Second)
	if err != nil {
		return nil, nil, err
	}
	slog.Info("Initialized DB connection")

	slog.Info("Seeding database")

	if err := seed.SeedToDb(db); err != nil {
		return nil, nil, err
	}

	slog.Info("Seeded Database")

	slog.Info("Initializing routers")

	serverConfig := api.NewAPIServer("0.0.0.0:"+env.PORT, db)
	app := serverConfig.InitaliseHTTPServer()

	slog.Info("Initialized routers")

	return app, func() {
		storage.CloseSqlite3(db)
	}, nil
}
