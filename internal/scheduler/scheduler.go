package scheduler

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/kondohiroki/go-boilerplate/config"
)

var UTC = time.UTC
var AsiaBangkok = time.FixedZone("Asia/Bangkok", 7*60*60)

func Start() {
	s := gocron.NewScheduler(UTC)
	s.SingletonModeAll()

	for _, schedule := range config.GetConfig().Schedules {
		if schedule.IsEnabled {
			switch schedule.Job {
			case "DoSomeThing":
				cronjob, err := s.CronWithSeconds(schedule.Cron).Do(func() {
					// j := job.NewJobContext()
					// j.DoSomeThing()
				})

				if err != nil {
					fmt.Printf("Failed to schedule SyncAll job: %v", err)
					continue
				}

				// Set up event listeners
				cronjob.SetEventListeners(func() {
					fmt.Println("DoSomeThing Job started -- round: ", cronjob.RunCount())
				}, func() {
					time.Sleep(1 * time.Second)

					// Print next run time in both utc and asia/bangkok
					fmt.Printf("\nNext run: %s / %s\n", cronjob.NextRun().UTC().String(), cronjob.NextRun().In(AsiaBangkok).String())

				})
			}
		}
	}

	fmt.Printf("Total jobs: %d jobs scheduled to run\n", len(s.Jobs()))
	fmt.Printf("Location: %s\n", s.Location().String())
	fmt.Println("Starting scheduler... (press Ctrl+C to quit)")

	s.StartImmediately()
	s.StartBlocking()
}
