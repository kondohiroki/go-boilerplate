package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"
	"github.com/kondohiroki/go-boilerplate/internal/helper/queue"
	"github.com/kondohiroki/go-boilerplate/internal/job"
	"github.com/kondohiroki/go-boilerplate/internal/logger"
	"github.com/kondohiroki/go-boilerplate/internal/repository"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddGroup(&cobra.Group{ID: "queue", Title: "Queue:"})
	rootCmd.AddCommand(
		queueWorkCommand,
		queueClearCommand,
		queueFlushCommand,
		queueForgetCommand,
		queueRetryCommand,
		queueRestoreCommand,
	)

	queueWorkCommand.Flags().StringP("queue", "q", "default", "(optional) queue name. for example: -q emails")
	queueWorkCommand.Flags().IntP("worker", "w", 1, "(optional) The number of worker goroutines to run. for example: -w 2")
	queueWorkCommand.Example = "  queue:work"
	queueWorkCommand.Example += "\n  queue:work -w 2"
	queueWorkCommand.Example += "\n  queue:work -q emails -w 2"

	queueRetryCommand.Flags().StringP("queue", "q", "default", "(optional) queue name. for example: -q emails")
	queueRetryCommand.Flags().StringP("id", "i", "", "(optional) job id. for example: --id df6df3af-d53d-49c2-bd50-80ba1d32b17b")
	queueRetryCommand.Example = `  queue:retry -q emails -i df6df3af-d53d-49c2-bd50-80ba1d32b17b`

	queueClearCommand.Flags().StringP("queue", "q", "default", "(optional) queue name. for example: -q emails")
	queueClearCommand.Flags().BoolP("all", "a", false, "(optional) force delete all jobs for all queues.")
	queueClearCommand.Example = "  queue:clear"
	queueClearCommand.Example += "\n  queue:clear -q emails"
	queueClearCommand.Example += "\n  queue:clear -a"

	queueFlushCommand.Flags().StringP("queue", "q", "default", "(optional) queue name. for example: -q emails")
	queueFlushCommand.Flags().BoolP("all", "a", false, "(optional) force delete all failed_jobs for all queues.")
	queueFlushCommand.Example = "  queue:flush"
	queueFlushCommand.Example += "\n  queue:flush -q emails"
	queueFlushCommand.Example += "\n  queue:flush -a"

	queueForgetCommand.Flags().StringP("id", "i", "", "(optional) job id. for example: --id df6df3af-d53d-49c2-bd50-80ba1d32b17b")
	queueForgetCommand.MarkFlagRequired("id")
	queueForgetCommand.Example = `  queue:forget -i df6df3af-d53d-49c2-bd50-80ba1d32b17b`

	queueRestoreCommand.Flags().StringP("queue", "q", "default", "(optional) queue name. for example: -q emails")
	queueRestoreCommand.Example = "  queue:restore"
	queueRestoreCommand.Example += "\n  queue:restore -q emails"
}

var queueWorkCommand = &cobra.Command{
	Use:     "queue:work",
	Short:   "Listen to a given queue",
	GroupID: "queue",
	Run: func(cmd *cobra.Command, _ []string) {
		// Setup all the required dependencies
		setupAll()

		queueName, _ := cmd.Flags().GetString("queue")
		numberOfWorkers, _ := cmd.Flags().GetInt("worker")

		logger.Log.Info(fmt.Sprintf("Starting %d queue workers for queue %s", numberOfWorkers, queueName))

		// Create a context that gets canceled when the program receives a termination signal.
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			sig := <-sigCh
			logger.Log.Info(fmt.Sprintf("Received signal: %v, canceling context", sig))
			cancel()
		}()

		q := queue.NewQueue(queueName)

		var wg sync.WaitGroup
		wg.Add(numberOfWorkers)

		for i := 0; i < numberOfWorkers; i++ {
			go func() {
				defer wg.Done()
				err := q.Run(ctx)
				if err != nil && err != context.Canceled {
					logger.Log.Error("Queue worker stopped with error", zap.Error(err))
				} else {
					logger.Log.Info("Queue worker stopped gracefully")
				}
			}()
		}

		// Wait for all workers to finish before exiting the program
		wg.Wait()

	},
}

var queueRetryCommand = &cobra.Command{
	Use:     "queue:retry",
	Short:   "Retry a failed queue job",
	GroupID: "queue",
	Run: func(cmd *cobra.Command, _ []string) {
		ctx := cmd.Context()

		// Setup all the required dependencies
		setupAll()

		queueName, _ := cmd.Flags().GetString("queue")
		jobID, _ := cmd.Flags().GetString("id")

		q := queue.NewQueue(queueName)

		if jobID != "" {
			uuid, err := uuid.Parse(jobID)
			if err != nil {
				logger.Log.Error("Cannot parse job id", zap.Error(err))
				return
			}

			err = q.RetryFailedByJobID(ctx, uuid)
			if err != nil {
				logger.Log.Error("Queue retry failed", zap.String("job_id", jobID), zap.Error(err))
			} else {
				logger.Log.Info(fmt.Sprintf("Queue retry completed. Job %s retried", jobID))
			}

		} else { // retry all failed jobs
			totalFailed, err := q.RetryAllFailed(ctx)
			if err != nil {
				logger.Log.Error("Queue retry failed", zap.Error(err))
			} else {
				logger.Log.Info(fmt.Sprintf("Queue retry completed. %d jobs retried", totalFailed))
			}
		}

	},
}

var queueClearCommand = &cobra.Command{
	Use:     "queue:clear",
	Short:   "Delete all of the jobs from the specified queue",
	GroupID: "queue",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, _ []string) {
		ctx := cmd.Context()
		// Setup all the required dependencies
		setupAll()

		queueName, _ := cmd.Flags().GetString("queue")
		all, _ := cmd.Flags().GetBool("all")

		if all {
			totalDeleted, err := queue.ClearAll(ctx)
			if err != nil {
				logger.Log.Error("Queue clear failed", zap.Error(err))
			} else {
				logger.Log.Info(fmt.Sprintf("Queue clear completed. %d queues deleted", totalDeleted))
			}
			return
		}

		q := queue.NewQueue(queueName)
		totalDeleted, err := q.Clear(ctx)
		if err != nil {
			logger.Log.Error("Queue clear failed", zap.Error(err))
		} else {
			logger.Log.Info(fmt.Sprintf("Queue clear completed. Queue %d deleted", totalDeleted))
		}

	},
}

var queueFlushCommand = &cobra.Command{
	Use:     "queue:flush",
	Short:   "Flush all of the failed queue jobs",
	GroupID: "queue",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, _ []string) {
		ctx := cmd.Context()

		// Setup all the required dependencies
		setupAll()

		queueName, _ := cmd.Flags().GetString("queue")
		all, _ := cmd.Flags().GetBool("all")

		if all {
			totalDeleted, err := queue.FlushAllFailed(ctx)
			if err != nil {
				logger.Log.Error("Queue clear failed_jobs failed", zap.Error(err))
			} else {
				logger.Log.Info(fmt.Sprintf("Queue clear failed_jobs completed. %d queues deleted", totalDeleted))
			}
			return
		}

		q := queue.NewQueue(queueName)
		totalDeleted, err := q.RemoveAllFailed(ctx)
		if err != nil {
			logger.Log.Error("Queue clear failed_jobs failed", zap.Error(err))
		} else {
			logger.Log.Info(fmt.Sprintf("Queue clear failed_jobs completed. Queue %d deleted", totalDeleted))
		}

	},
}

var queueForgetCommand = &cobra.Command{
	Use:     "queue:forget",
	Short:   "Delete a failed queue job",
	GroupID: "queue",
	Run: func(cmd *cobra.Command, _ []string) {
		ctx := cmd.Context()

		// Setup all the required dependencies
		setupAll()

		jobID, _ := cmd.Flags().GetString("id")

		if jobID != "" {
			uuid, err := uuid.Parse(jobID)
			if err != nil {
				logger.Log.Error("Cannot parse job id", zap.Error(err))
				return
			}

			queueName, err := queue.RemoveJobOnAnyQueueByID(ctx, uuid)
			if err != nil {
				logger.Log.Error("Queue forget failed", zap.Error(err))
			} else {
				logger.Log.Info(fmt.Sprintf("Queue forget completed. Job %s deleted from queue %s", jobID, queueName))
			}
		}

	},
}

var queueRestoreCommand = &cobra.Command{
	Use:     "queue:restore",
	Short:   "Restore a failed and unfinished job to the redis queue",
	GroupID: "queue",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		// Setup all the required dependencies
		setupAll()

		queueName, _ := cmd.Flags().GetString("queue")
		q := queue.NewQueue(queueName)

		repo := repository.NewRepository()

		// Get all pending jobs and restore them to the redis queue
		unfinishedJobs, err := repo.Job.GetUnfinishedJobs(ctx)
		if err != nil {
			logger.Log.Error("Get pending jobs error", zap.Error(err))
			return
		}

		if len(unfinishedJobs) == 0 {
			logger.Log.Info("No unfinished jobs found")
			return
		}

		if len(unfinishedJobs) > 0 {
			logger.Log.Info(fmt.Sprintf("Found %d unfinished jobs", len(unfinishedJobs)))

			// reset all processing jobs to pending
			err = repo.Job.ResetProcessingJobsToPending(ctx)
		}

		for _, j := range unfinishedJobs {
			jobItem, err := job.NewJob(
				j.HandlerName,
				j.Payload,
				j.MaxAttempts,
				j.Delay,
			)
			if err != nil {
				logger.Log.Error("New job error", zap.Error(err))
				continue
			}

			if j.Status == job.StatusFailed {
				err = q.EnqueueFailedJobs(ctx, jobItem)
			} else if j.Status == job.StatusPending {
				err = q.EnqueuePendingJobs(ctx, jobItem)
			}

			if err != nil {
				logger.Log.Error("Restore job error", zap.Error(err))
			} else {
				logger.Log.Info(fmt.Sprintf("Queue restore completed. Job %s restored to queue %s", j.ID, queueName))
			}
		}

	},
}
