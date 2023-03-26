package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/kondohiroki/go-boilerplate/config"
	"github.com/kondohiroki/go-boilerplate/internal/db"
	"github.com/kondohiroki/go-boilerplate/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const defaultConfigFile = "config/config.yaml"

var configFile string
var rootCmd = &cobra.Command{
	Use: func() string {
		if nameForCLI := viper.GetString("app.nameForCLI"); nameForCLI != "" {
			return nameForCLI
		}
		return "my-app"
	}(),
	Short: "Made with ❤️",
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Usage()
	},
	Version: "0.0.1",
}

func init() {
	cobra.OnInitialize(
		setUpConfig,
		setUpLogger,
		setUpPostgres,
		setUpSentry,
	)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", fmt.Sprintf("config file (default is %s)", defaultConfigFile))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("rootCmd.Execute() Error: %v", err)
		os.Exit(1)
	}
}

func setUpConfig() {
	if configFile == "" {
		configFile = defaultConfigFile
	}

	config.SetConfig(configFile)
}

func setUpLogger() {
	logger.InitLogger("zap")
}

func setUpPostgres() {
	// Create the database connection pool
	if config.GetConfig().Postgres.Host != "" {
		if config.GetConfig().Postgres.Schema == "" {
			log.Fatalf("Postgres schema is not set")
		}

		// Initialize database schema if it doesn't exist
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		err := db.InitSchema(ctx, config.GetConfig().Postgres, config.GetConfig().Postgres.Schema)
		if err != nil {
			log.Fatalf("db.InitSchema() Error: %v", err)
		}

		err = db.InitPgConnectionPool(config.GetConfig().Postgres)
		if err != nil {
			log.Fatalf("db.InitPgConnectionPool() Error: %v", err)
		}
	}

}

func setUpSentry() {
	// Don't initialize sentry if DSN is not set
	if config.GetConfig().Sentry.Dsn == "" {
		return
	}

	// Initialize sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn: config.GetConfig().Sentry.Dsn,
		// BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
		// 	if hint.Context != nil {
		// 		if c, ok := hint.Context.Value(sentry.RequestContextKey).(*fiber.Ctx); ok {
		// 			// You have access to the original Context if it panicked
		// 			fmt.Println(utils.CopyString(c.Hostname()))
		// 		}
		// 	}
		// 	fmt.Println(event)
		// 	return event
		// },
		Debug:            config.GetConfig().Sentry.Debug,
		AttachStacktrace: true,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		Environment:      config.GetConfig().Sentry.Environment,
		Release:          config.GetConfig().Sentry.Release,
	})

	if err != nil {
		logger.Log.Error("Creata Sentry instant error: %v", zap.Error(err))
	} else {
		logger.Log.Info("Creata Sentry instant success")
	}

	defer sentry.Flush(2 * time.Second)
}
