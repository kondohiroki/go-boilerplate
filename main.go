package main

import (
	"log"
	"os"

	"github.com/kondohiroki/go-boilerplate/cmd"
	"github.com/kondohiroki/go-boilerplate/internal/db/pgx"
	"github.com/kondohiroki/go-boilerplate/internal/logger"
)

func main() {
	defer func() {
		log.Default().Println("Cleaning up...")

		log.Default().Println("Removing Liveness file /tmp/live")
		_ = os.Remove("/tmp/live")

		// Close the database connection pool
		log.Default().Println("Closing database connection pool")
		pgx.ClosePgxPool()

		// Flush the log buffer
		log.Default().Println("Flushing the log buffer")
		logger.Log.Sync()
	}()

	// Liveness probe for Kubernetes
	log.Default().Println("Creating Liveness file /tmp/live")
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatalf("Cannot create a Liveness file: %v", err)
	}

	// Start the app here via CLI commands
	cmd.Execute()
}
