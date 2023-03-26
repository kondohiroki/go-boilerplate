package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kondohiroki/go-boilerplate/config"
	"github.com/kondohiroki/go-boilerplate/internal/scheduler"
	"github.com/lnquy/cron"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(
		listScheduleCommand,
		startScheduleCommand,
	)
}

var startScheduleCommand = &cobra.Command{
	Use:   "schedule:run",
	Short: "Start schedule job",
	Run: func(_ *cobra.Command, _ []string) {
		printScheduleList()
		scheduler.Start()
	},
}

var listScheduleCommand = &cobra.Command{
	Use:   "schedule:list",
	Short: "List all schedule jobs",
	Run: func(_ *cobra.Command, _ []string) {

		printScheduleList()

	},
}

func printScheduleList() {
	exprDesc, _ := cron.NewDescriptor()

	// Print the job list as a table in the console
	tableWriter := table.NewWriter()
	tableWriter.SetOutputMirror(os.Stdout)
	tableWriter.AppendHeader(table.Row{"No.", "Job Name", "Cron Expression", "Schedule"})
	for i, schedule := range config.GetConfig().Schedules {
		desc, _ := exprDesc.ToDescription(schedule.Cron, cron.Locale_en)

		tableWriter.AppendRow(table.Row{
			i + 1,
			schedule.Job,
			schedule.Cron,
			desc,
		})

	}

	tableWriter.Render()
}
