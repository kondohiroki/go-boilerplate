package cmd

import (
	"github.com/kondohiroki/go-boilerplate/internal/app"
	"github.com/kondohiroki/go-boilerplate/internal/jobs"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(
		notifyDiscordCommand(),
	)
}

func notifyDiscordCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "notify:discord",
		Short: "Notify to Discord",
		Run: func(cmd *cobra.Command, args []string) {
			jobs.SendNotificationViaDiscord(app.GetAppContext())
		},
	}
}
