package main

import (
	"log"
	"os"

	"github.com/kondohiroki/go-boilerplate/cmd"
	"github.com/kondohiroki/go-boilerplate/internal/app"
)

func main() {
	defer func() {
		_ = os.Remove("/tmp/live")

		// Close Discord connection
		if app.GetAppContext() != nil && app.GetAppContext().Discord != nil {
			_ = app.GetAppContext().Discord.Close()
		}
	}()

	// Liveness probe for Kubernetes
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatalf("Cannot create a Liveness file: %v", err)
	}

	// Start the app here via CLI commands
	cmd.Execute()
}
