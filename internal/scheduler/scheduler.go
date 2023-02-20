package scheduler

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/kondohiroki/go-boilerplate/internal/app"
)

func Start(c *app.AppContext) {
	s := gocron.NewScheduler(time.UTC)

	fmt.Printf("Starting scheduler with %d schedules\n", len(c.Config.Schedules))

	for _, schedule := range c.Config.Schedules {
		if schedule.IsEnabled {
			switch schedule.Job {
			case "OpenMergeRequestToSIT":
				// s.CronWithSeconds(schedule.Cron).Do(func() { jobs.OpenMergeRequestToSIT(appCtx) })
			}
		}
	}

	s.StartBlocking()
}
