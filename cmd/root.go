package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kondohiroki/go-boilerplate/config"
	"github.com/kondohiroki/go-boilerplate/internal/app"
	"github.com/kondohiroki/go-boilerplate/internal/discord"
	"github.com/kondohiroki/go-boilerplate/internal/logger"
	"github.com/spf13/cobra"
)

const defaultConfigFile = "config/config.yaml"

var configFile string
var rootCmd = &cobra.Command{
	Use:   "myapp",
	Short: "made with ❤️",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
	Version: "0.0.1",
}

func init() {
	cobra.OnInitialize(
		setConfigFile,
		setUpAppContext,
	)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", fmt.Sprintf("config file (default is %s)", defaultConfigFile))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("rootCmd.Execute() Error: %v", err)
		os.Exit(1)
	}
}

func setConfigFile() {
	if configFile == "" {
		configFile = defaultConfigFile
	}
}

func setUpAppContext() {
	// Initialize AppContext
	appCtx := &app.AppContext{
		Ctx: context.Background(),
	}

	config.SetConfig(configFile)
	appCtx.Config = config.GetConfig()

	// Create logger
	logger := logger.NewLogger()
	appCtx.Logger = logger

	// Create Discord client
	if appCtx.Config.Discord.Token != "" {
		appCtx.Discord = discord.NewDiscord(appCtx.Config.Discord.Token)
	}

	app.SetAppContext(appCtx)
}
