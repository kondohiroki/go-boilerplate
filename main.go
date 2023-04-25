package main

import (
	"log"
	"os"

	"github.com/kondohiroki/go-boilerplate/cmd"
	"github.com/kondohiroki/go-boilerplate/internal/db/pgx"
	"github.com/kondohiroki/go-boilerplate/internal/logger"

	"go.uber.org/automaxprocs/maxprocs"
)

func main() {
	defer func() {
		_ = os.Remove("/tmp/live")

		// Close the database connection pool
		pgx.ClosePgxPool()

		// Flush the log buffer
		if logger.Log != nil {
			logger.Log.Sync()
		}
	}()

	// Liveness probe for Kubernetes
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatalf("Cannot create a Liveness file: %v", err)
	}

	// Set the maximum number of CPUs to use
	nopLog := func(string, ...interface{}) {}
	_, err = maxprocs.Set(maxprocs.Logger(nopLog))
	if err != nil {
		log.Fatalf("Cannot set maxprocs: %v", err)
	}

	// Start the app here via CLI commands
	cmd.Execute()
}
