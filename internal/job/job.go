package job

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/kondohiroki/go-boilerplate/internal/repository"
)

type JobStatus string

const (
	JOB_PENDING         JobStatus = "PENDING"
	JOB_IN_PROGRESS     JobStatus = "IN PROGRESS"
	JOB_SUCCESS         JobStatus = "SUCCESS"
	JOB_FAILED          JobStatus = "FAILED"
	JOB_NOTHING_CHANGED JobStatus = "NOTHING CHANGED"
)

type JobFunc func(ctx context.Context) (JobStatus, time.Time, int64, error)

type JobChain struct {
	JobName string
	JobFunc JobFunc
}

type JobContext struct {
	Ctx  context.Context
	Repo *repository.Repository
}

type JobReport struct {
	JobName       string
	Status        JobStatus
	Error         error
	StartedTime   time.Time
	ExecutionTime int64
}

func NewJobContext() *JobContext {
	return &JobContext{
		Ctx:  context.Background(),
		Repo: repository.NewRepository(),
	}
}

func PrintJobReport(jobReports []JobReport) {
	// Define the color for each status value
	var statusColors = map[JobStatus]text.Colors{
		JOB_PENDING:         {text.FgBlue},
		JOB_IN_PROGRESS:     {text.FgYellow},
		JOB_SUCCESS:         {text.FgGreen},
		JOB_FAILED:          {text.FgRed},
		JOB_NOTHING_CHANGED: {text.FgGreen},
	}

	tableWriter := table.NewWriter()
	tableWriter.SetOutputMirror(os.Stdout)
	tableWriter.AppendHeader(table.Row{"No.", "Job Name", "Status", "Error", "Started Time", "Execution Time"})
	tableWriter.SetColumnConfigs([]table.ColumnConfig{
		{Name: "No.", AlignHeader: text.AlignCenter, Align: text.AlignLeft},
		{Name: "Job Name", AlignHeader: text.AlignCenter, Align: text.AlignLeft},
		{Name: "Status", AlignHeader: text.AlignCenter, Align: text.AlignLeft},
		{Name: "Error", AlignHeader: text.AlignCenter, Align: text.AlignLeft},
		{Name: "Started Time", AlignHeader: text.AlignCenter, Align: text.AlignLeft},
		{Name: "Execution Time", AlignHeader: text.AlignCenter, Align: text.AlignLeft},
	})

	for index, jobReport := range jobReports {
		rowColor := statusColors[jobReport.Status]
		var errorMsg any
		if jobReport.Error != nil {
			errorMsg = jobReport.Error
		} else {
			errorMsg = "-"
		}

		row := table.Row{
			strconv.Itoa(index + 1),
			jobReport.JobName,
			rowColor.Sprint(jobReport.Status),
			errorMsg,
			jobReport.StartedTime.Format(time.RFC822),
			fmt.Sprintf("%dms", jobReport.ExecutionTime),
		}
		tableWriter.AppendRow(row)
	}

	tableWriter.Render()
}
