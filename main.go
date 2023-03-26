package main

import (
	"log"
	"os"

	"github.com/kondohiroki/go-boilerplate/cmd"
	"github.com/kondohiroki/go-boilerplate/internal/db"
	"github.com/kondohiroki/go-boilerplate/internal/logger"
)

func main() {
	defer func() {
		_ = os.Remove("/tmp/live")

		// Close the database connection pool
		db.ClosePgxPool()

		// Flush the log buffer
		logger.Log.Sync()
	}()

	// Liveness probe for Kubernetes
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatalf("Cannot create a Liveness file: %v", err)
	}

	// Start the app here via CLI commands
	cmd.Execute()
}
