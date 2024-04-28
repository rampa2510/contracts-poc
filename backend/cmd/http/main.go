package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rampa2510/contracts-poc/config"
	"github.com/rampa2510/contracts-poc/internal/storage"
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
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}

	cleanup, err := run(env)

	defer cleanup()
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}
}

func run(env config.EnvVars) (func(), error) {
	cleanup, err := buildServer(env)
	if err != nil {
		return nil, err
	}

	// start the server
	// go func() {
	// 	app.Listen("0.0.0.0:" + env.PORT)
	// }()
	//
	// return a function to close the server and database
	return func() {
		cleanup()
		// app.Shutdown()
	}, nil
}

func buildServer(env config.EnvVars) (func(), error) {
	fmt.Println("Initializing DB connection")
	// init the storage
	db, err := storage.BootstrapSqlite3(env.DB_FILE, 10*time.Second)
	if err != nil {
		return nil, err
	}
	fmt.Println("Initialized DB connection")

	return func() {
		storage.CloseSqlite3(db)
	}, nil
}
