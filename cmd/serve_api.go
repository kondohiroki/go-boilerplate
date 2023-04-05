package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/kondohiroki/go-boilerplate/config"
	http_api "github.com/kondohiroki/go-boilerplate/internal/interface/httpapi"
	"github.com/kondohiroki/go-boilerplate/internal/interface/httpapi/validation"
	"github.com/kondohiroki/go-boilerplate/internal/logger"
	"github.com/kondohiroki/go-boilerplate/internal/router"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveAPICmd)
}

var serveAPICmd = &cobra.Command{
	Use:   "serve-api",
	Short: "Start the RESTful API",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create http router
		r := router.NewFiberRouter()

		// Register routes
		http_api.RegisterRoute(r)

		// Create validator instance
		validation.InitValidator()

		// Get port from config
		port := config.GetConfig().HttpServer.Port

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()
		go func() {
			localIP, _ := getLocalIP()

			logger.Log.Info(fmt.Sprintf("Starting server on port %d", port))
			logger.Log.Info(fmt.Sprintf("Local: http://localhost:%d", port))
			logger.Log.Info(fmt.Sprintf("Network: http://%s:%d", localIP, port))

			if err := r.Listen(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
				logger.Log.Fatal(fmt.Sprintf("listen: %s\n", err))
			}
		}()

		<-ctx.Done()
		stop()
		fmt.Println("\nShutting down gracefully, press Ctrl+C again to force")

		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := r.ShutdownWithTimeout(5 * time.Second); err != nil {
			fmt.Println(err)
		}

		return nil
	},
}

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("local IP not found")
}
