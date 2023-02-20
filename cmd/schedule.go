package cmd

import (
	"fmt"

	"github.com/kondohiroki/go-boilerplate/internal/app"
	"github.com/kondohiroki/go-boilerplate/internal/scheduler"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(
		listScheduleCommand(),
		startScheduleCommand(),
	)
}

func listScheduleCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "schedule:list",
		Short: "List all schedule jobs",
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Printf("%-20s%-15s%s\n", "Cron", "IsEnabled", "Job")
			for _, schedule := range app.GetAppContext().Config.Schedules {
				fmt.Printf("%-20s%-15t%s\n", schedule.Cron, schedule.IsEnabled, schedule.Job)
			}

		},
	}
}

func startScheduleCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "schedule:run",
		Short: "Start schedule job",
		Run: func(cmd *cobra.Command, args []string) {
			scheduler.Start(app.GetAppContext())
		},
	}
}
