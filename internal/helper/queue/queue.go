package queue

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"github.com/kondohiroki/go-boilerplate/internal/db/model"
	"github.com/kondohiroki/go-boilerplate/internal/db/rdb"
	"github.com/kondohiroki/go-boilerplate/internal/job"
	"github.com/kondohiroki/go-boilerplate/internal/logger"
	"github.com/kondohiroki/go-boilerplate/internal/repository"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

/*
Queue is a FIFO.
Implement redis BLMOVE with the RIGHT and LEFT arguments.
*/

type Queue struct {
	Key              string
	KeyWithoutPrefix string

	repo *repository.Repository
}

// QueueInfo holds information about a specific queue.
type QueueInfo struct {
	Key              string `json:"key"`
	KeyWithoutPrefix string `json:"key_without_prefix"`
	NumberOfItems    int64  `json:"number_of_items"`
}

func NewQueue(key string) *Queue {
	return &Queue{
		Key:              rdb.AddQueuePrefix(key),
		KeyWithoutPrefix: key,
		repo:             repository.NewRepository(),
	}
}

// Adds an item to the source list (the end of the queue).
func (q *Queue) Enqueue(ctx context.Context, jobs ...*job.Job) error {
	rdbClient := rdb.GetRedisClient()

	for _, j := range jobs {
		jobBytes, err := sonic.Marshal(j)
		if err != nil {
			return err
		}

		// Add job to postgres for backup
		_, err = q.repo.Job.AddJob(ctx, model.Job{
			ID:          j.ID,
			Queue:       q.KeyWithoutPrefix,
			HandlerName: j.HandlerName,
			Payload:     j.Payload,
			MaxAttempts: j.MaxAttempts,
			Delay:       j.Delay,
			Status:      job.StatusPending,
			CreatedAt:   j.CreatedAt,
		})
		if err != nil {
			logger.Log.Error("Error adding job to postgres", zap.Error(err))
			return err
		}

		// Add job to redis
		err = rdbClient.LPush(ctx, q.Key, jobBytes).Err()
		if err != nil {
			return err
		}

		if err != nil {
			logger.Log.Error("Error adding job to redis", zap.Error(err))
			return err
		}
	}

	return nil
}

// Restore pending jobs from postgres to redis.
func (q *Queue) EnqueuePendingJobs(ctx context.Context, jobs ...*job.Job) error {
	rdbClient := rdb.GetRedisClient()

	for _, j := range jobs {
		jobBytes, err := sonic.Marshal(j)
		if err != nil {
			return err
		}
		if err != nil {
			logger.Log.Error("Error adding job to postgres", zap.Error(err))
			return err
		}

		// Add job to redis
		err = rdbClient.LPush(ctx, q.Key, jobBytes).Err()
		if err != nil {
			return err
		}

		if err != nil {
			logger.Log.Error("Error adding job to redis", zap.Error(err))
			return err
		}
	}

	return nil
}

// Restore failed jobs from postgres to redis.
func (q *Queue) EnqueueFailedJobs(ctx context.Context, jobs ...*job.Job) error {
	rdbClient := rdb.GetRedisClient()

	for _, j := range jobs {
		jobBytes, err := sonic.Marshal(j)
		if err != nil {
			return err
		}

		if err != nil {
			logger.Log.Error("Error adding job to postgres", zap.Error(err))
			return err
		}

		// Add job to redis
		err = rdbClient.LPush(ctx, q.Key+"_failed", jobBytes).Err()
		if err != nil {
			return err
		}

		if err != nil {
			logger.Log.Error("Error adding job to redis", zap.Error(err))
			return err
		}
	}

	return nil
}

// Removes an item from the source list (the start of the queue) and adds it to the destkey list (temporary storage location).
func (q *Queue) Dequeue(ctx context.Context, timeout time.Duration) (*job.Job, error) {
	sourceKey := q.Key
	destKey := q.Key + "_attempt"
	rdbClient := rdb.GetRedisClient()

	// Move the job from the source list to the temporary list
	result, err := rdbClient.BLMove(ctx, sourceKey, destKey, "RIGHT", "LEFT", timeout).Result()

	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var j job.Job
	err = sonic.Unmarshal([]byte(result), &j)
	if err != nil {
		return nil, err
	}

	// Update status of job in postgres
	if err := q.repo.Job.UpdateJobStatus(ctx, j.ID, job.StatusProcessing); err != nil {
		logger.Log.Error("Error updating job status in postgres", zap.Error(err))
		return nil, err
	}

	// Check if the job has reached the maximum number of attempts
	if j.MaxAttempts > 0 && j.Attempts >= j.MaxAttempts {
		// Remove the job from the temporary list, it has reached the maximum number of attempts
		_, _ = rdbClient.LRem(ctx, destKey, 1, result).Result()
		return nil, fmt.Errorf("job %s reached the maximum number of attempts (%d)", j.ID, j.MaxAttempts)
	}

	// Increment the number of attempts and update the job in the temporary list
	j.Attempts++
	jobBytes, err := sonic.Marshal(&j)
	if err != nil {
		return nil, err
	}

	// Remove the job from the temporary list
	_, err = rdbClient.LRem(ctx, destKey, 1, result).Result()
	if err != nil {
		return nil, err
	}

	// Add the updated job to the temporary list
	err = rdbClient.LPush(ctx, destKey, jobBytes).Err()
	if err != nil {
		return nil, err
	}

	return &j, nil
}

// Removes the processed item with the given job ID from the destkey list (temporary storage location).
func (q *Queue) RemoveProcessed(ctx context.Context, jobID uuid.UUID, jobError error) error {
	destKey := q.Key + "_attempt"
	sourceKey := q.Key
	failedJobsKey := q.Key + "_failed"
	rdbClient := rdb.GetRedisClient()

	queues := rdbClient.LRange(ctx, destKey, 0, -1)
	for _, queueItem := range queues.Val() {
		var j job.Job
		if err := sonic.Unmarshal([]byte(queueItem), &j); err != nil {
			return err
		}

		if j.ID == jobID {
			// Remove the job from the temporary list
			if err := removeJobFromList(ctx, rdbClient, destKey, queueItem); err != nil {
				return err
			}

			// if the job failed, add it to the failed_jobs list
			if jobError != nil {
				j.Errors = append(j.Errors, jobError.Error())
				return handleFailedJob(ctx, rdbClient, q.KeyWithoutPrefix, q.repo, j, sourceKey, failedJobsKey)
			}

			// if the job was successful, then update the job status to completed in postgres
			if err := q.repo.Job.UpdateJobStatus(ctx, j.ID, job.StatusCompleted); err != nil {
				logger.Log.Error("Error updating job status in postgres", zap.Error(err))
				return err
			}

			break
		}
	}

	return nil
}

// Remove the job from the temporary list if it was processed successfully
func removeJobFromList(ctx context.Context, rdbClient redis.Cmdable, destKey string, queueItem string) error {
	// Remove the job from the temporary list
	if _, err := rdbClient.LRem(ctx, destKey, 1, queueItem).Result(); err != nil {
		return err
	}
	return nil
}

// If the job failed, add it to the failed_jobs list
func handleFailedJob(ctx context.Context, rdbClient redis.Cmdable, queue string, repo *repository.Repository, j job.Job, sourceKey string, failedJobsKey string) error {
	if j.MaxAttempts == 0 || j.Attempts < j.MaxAttempts {
		time.Sleep(time.Duration(j.Delay) * time.Second)
		jobBytes, err := sonic.Marshal(&j)
		if err != nil {
			return err
		}

		// Update job status in postgres
		if err := repo.Job.UpdateJobStatus(ctx, j.ID, job.StatusPending); err != nil {
			logger.Log.Error("Error updating job status in postgres", zap.Error(err))
			return err
		}

		// Add the job back to the source list (the beginning of the queue) with the delay
		err = rdbClient.LPush(ctx, sourceKey, jobBytes).Err()
		if err != nil {
			return err
		}
	} else {
		// update job status in postgres
		if err := repo.Job.UpdateJobStatus(ctx, j.ID, job.StatusFailed); err != nil {
			logger.Log.Error("Error updating job status in postgres", zap.Error(err))
			return err
		}

		// Add failed job to postgres
		_, err := repo.Job.AddFailedJob(ctx, model.FaildJob{
			JobID:    j.ID,
			Queue:    queue,
			Payload:  j.Payload,
			Error:    strings.Join(j.Errors, ","),
			FailedAt: time.Now(),
		})
		if err != nil {
			logger.Log.Error("Error adding failed job to postgres", zap.Error(err))
			return err
		}

		if err := addJobToFailedList(ctx, rdbClient, j, failedJobsKey); err != nil {
			return err
		}
	}
	return nil
}

func addJobToFailedList(ctx context.Context, rdbClient redis.Cmdable, job job.Job, failedJobsKey string) error {
	logger.Log.Info("Job has reached the maximum number of attempts. It will be added to the failed_jobs list", zap.String("job_id", job.ID.String()))
	jobBytes, err := sonic.Marshal(&job)
	if err != nil {
		return err
	}
	err = rdbClient.LPush(ctx, failedJobsKey, jobBytes).Err()
	if err != nil {
		return err
	}
	return nil
}

// RetryFailedByJobID moves a failed item with the given job ID from the destkey list back to the source list for retrying.
func (q *Queue) RetryFailedByJobID(ctx context.Context, jobID uuid.UUID) error {
	failedJobsKey := q.Key + "_failed"
	destkey := q.Key
	rdbClient := rdb.GetRedisClient()

	queues, err := rdbClient.LRange(ctx, failedJobsKey, 0, -1).Result()
	if err != nil {
		return err
	}

	found, err := findAndRetryJob(ctx, rdbClient, queues, jobID, failedJobsKey, destkey)
	if err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("job with ID %s not found in temporary list", jobID)
	}

	// update job status in postgres
	if err := q.repo.Job.UpdateJobStatus(ctx, jobID, job.StatusPending); err != nil {
		logger.Log.Error("Error updating job status in postgres", zap.Error(err))
		return err
	}

	// remove job from failed_jobs in postgres
	if err := q.repo.Job.RemoveFailedJob(ctx, jobID); err != nil {
		logger.Log.Error("Error removing job from failed_jobs in postgres", zap.Error(err))
		return err
	}

	return nil
}

func findAndRetryJob(ctx context.Context, rdbClient redis.Cmdable, queues []string, jobID uuid.UUID, failedJobsKey, destkey string) (bool, error) {
	for _, failedItem := range queues {
		var job job.Job
		if err := sonic.Unmarshal([]byte(failedItem), &job); err != nil {
			return false, err
		}

		if job.ID == jobID {
			if err := resetAndMoveJob(ctx, rdbClient, job, failedJobsKey, destkey, failedItem); err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}

func resetAndMoveJob(ctx context.Context, rdbClient redis.Cmdable, job job.Job, failedJobsKey, destkey, failedItem string) error {
	job.Attempts = 0
	updatedItem, err := sonic.Marshal(&job)
	if err != nil {
		return err
	}

	if _, err := rdbClient.LRem(ctx, failedJobsKey, 1, failedItem).Result(); err != nil {
		return err
	}

	err = rdbClient.RPush(ctx, destkey, updatedItem).Err()
	if err != nil {
		_ = rdbClient.LPush(ctx, failedJobsKey, failedItem).Err()
		return err
	}

	return nil
}

// Moves all failed items from the destkey list back to the source list for retrying.
func (q *Queue) RetryAllFailed(ctx context.Context) (int, error) {
	failedJobsKey := q.Key + "_failed"
	destkey := q.Key
	rdbClient := rdb.GetRedisClient()

	count := 0

	for {
		// Remove the failed item from the temporary list.
		failedItem, err := rdbClient.LPop(ctx, failedJobsKey).Result()
		if err == redis.Nil {
			// No more items left in the temporary list, break the loop.
			break
		}
		if err != nil {
			return count, err
		}

		// Reset the attempts counter
		var j job.Job
		if err := sonic.Unmarshal([]byte(failedItem), &j); err != nil {
			return count, err
		}
		j.Attempts = 0
		updatedItem, err := sonic.Marshal(&j)
		if err != nil {
			return count, err
		}

		// Add the failed item back to the source list with the reset attempts counter.
		err = rdbClient.RPush(ctx, destkey, updatedItem).Err()
		if err != nil {

			// Add the failed item back to the temporary list in case of an RPush error.
			_ = rdbClient.LPush(ctx, failedJobsKey, failedItem).Err()
			return count, err
		}

		// update job status in postgres
		if err := q.repo.Job.UpdateJobStatus(ctx, j.ID, job.StatusPending); err != nil {
			logger.Log.Error("Error updating job status in postgres", zap.Error(err))
			return count, err
		}

		// remove job from failed_jobs in postgres
		if err := q.repo.Job.RemoveFailedJob(ctx, j.ID); err != nil {
			logger.Log.Error("Error removing job from failed_jobs in postgres", zap.Error(err))
			return count, err
		}

		count++
	}

	return count, nil
}

// Returns the current length of the source list (the number of items in the queue).
func (q *Queue) Length(ctx context.Context) (int64, error) {
	rdbClient := rdb.GetRedisClient()

	length, err := rdbClient.LLen(ctx, q.Key).Result()
	if err != nil {
		return 0, err
	}

	return length, nil
}

// IsEmpty checks if the source list (queue) is empty.
func (q *Queue) IsEmpty(ctx context.Context) (bool, error) {
	rdbClient := rdb.GetRedisClient()

	length, err := rdbClient.LLen(ctx, q.Key).Result()
	if err != nil {
		return false, err
	}

	return length == 0, nil
}

// Clear removes all items from the source list (queue).
func (q *Queue) Clear(ctx context.Context) (int64, error) {
	rdbClient := rdb.GetRedisClient()

	// Get the length of the queue before deleting the key.
	length, err := rdbClient.LLen(ctx, q.Key).Result()
	if err != nil {
		return 0, fmt.Errorf("error getting length of key %s: %w", q.Key, err)
	}

	// Remove all items from the source list.
	_, err = rdbClient.Del(ctx, q.Key).Result()
	if err != nil {
		return 0, err
	}

	return length, nil
}

// RemoveJobByID removes the job with the matching job ID from the source list.
func (q *Queue) RemoveJobByID(ctx context.Context, jobID uuid.UUID) (bool, error) {
	rdbClient := rdb.GetRedisClient()

	queues, err := rdbClient.LRange(ctx, q.Key, 0, -1).Result()
	if err != nil {
		return false, err
	}

	for _, queueItem := range queues {
		var job job.Job
		if err := sonic.Unmarshal([]byte(queueItem), &job); err != nil {
			return false, err
		}

		if job.ID == jobID {
			// Remove the job with the matching ID from the source list.
			if _, err := rdbClient.LRem(ctx, q.Key, 1, queueItem).Result(); err != nil {
				return false, err
			}

			return true, nil
		}
	}

	return false, nil
}

// RemoveFailedByID removes the failed item with the matching job ID from the failed list.
func (q *Queue) RemoveFailedByID(ctx context.Context, jobID uuid.UUID) error {
	failedJobsKey := q.Key + "_failed"
	rdbClient := rdb.GetRedisClient()

	queues, err := rdbClient.LRange(ctx, failedJobsKey, 0, -1).Result()
	if err != nil {
		return err
	}

	found := false
	for _, failedItem := range queues {
		var job job.Job
		if err := sonic.Unmarshal([]byte(failedItem), &job); err != nil {
			return err
		}

		if job.ID == jobID {
			found = true

			// Remove the job with the matching ID from the failed list.
			if _, err := rdbClient.LRem(ctx, failedJobsKey, 1, failedItem).Result(); err != nil {
				return err
			}

			break
		}
	}

	if !found {
		return fmt.Errorf("job with ID %s not found in temporary list", jobID)
	}

	return nil
}

// RemoveAllFailed removes all items from the failed list.
func (q *Queue) RemoveAllFailed(ctx context.Context) (int64, error) {
	rdbClient := rdb.GetRedisClient()

	// Get the length of the queue before deleting the key.
	length, err := rdbClient.LLen(ctx, q.Key).Result()
	if err != nil {
		return 0, fmt.Errorf("error getting length of key %s: %w", q.Key, err)
	}

	// Remove all items from the failed list.
	_, err = rdbClient.Del(ctx, q.Key+"_failed").Result()
	if err != nil {
		return 0, err
	}

	return length, nil
}

// Peek returns the first N items in the source list without removing them.
func (q *Queue) Peek(ctx context.Context, count int64) ([]interface{}, error) {
	rdbClient := rdb.GetRedisClient()

	rawItems, err := rdbClient.LRange(ctx, q.Key, 0, count-1).Result()
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, len(rawItems))
	for i, rawItem := range rawItems {
		var item interface{}
		err = sonic.Unmarshal([]byte(rawItem), &item)
		if err != nil {
			return nil, err
		}
		items[i] = item
	}

	return items, nil
}

func (q *Queue) Run(ctx context.Context) error {
	handlerMap := job.NewHandlerMap()
	waitingMessagePrinted := false

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Context canceled, stopping the Queue")
			return ctx.Err()
		default:
			err, printed := processJob(ctx, q, handlerMap, waitingMessagePrinted)
			if err != nil {
				logger.Log.Error("Error processing job", zap.Error(err))
			}
			waitingMessagePrinted = printed
		}
	}
}

func processJob(ctx context.Context, q *Queue, handlerMap job.HandlerMap, waitingMessagePrinted bool) (err error, printed bool) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred while processing job: %v", r)
			printed = waitingMessagePrinted
			logger.Log.Error("Recovered from panic", zap.Any("panic", r))
		}
	}()

	dequeuedJob, err := q.Dequeue(ctx, time.Second*5)
	if err != nil {
		return fmt.Errorf("error dequeueing job: %w", err), waitingMessagePrinted
	}

	if dequeuedJob == nil {
		if !waitingMessagePrinted {
			logger.Log.Info(fmt.Sprintf("waiting for %s ...", q.KeyWithoutPrefix))
			waitingMessagePrinted = true
		}
		return nil, waitingMessagePrinted
	}

	waitingMessagePrinted = false

	logger.Log.Info("Starting job", zap.String("ID", dequeuedJob.ID.String()))

	handlerFunc, ok := handlerMap[dequeuedJob.HandlerName]
	if !ok {
		err := fmt.Errorf("handler not found: %v", dequeuedJob.HandlerName)
		q.RemoveProcessed(ctx, dequeuedJob.ID, err)
		return err, waitingMessagePrinted
	}

	handler := handlerFunc()
	err = sonic.Unmarshal(dequeuedJob.Payload, handler)
	if err != nil {
		return fmt.Errorf("error unmarshaling job payload: %w", err), waitingMessagePrinted
	}

	logger.Log.Info("Processing job", zap.String("ID", dequeuedJob.ID.String()), zap.String("handler", dequeuedJob.HandlerName))
	handlerError := handler.Handle()

	logger.Log.Info("Finished processing job", zap.String("ID", dequeuedJob.ID.String()), zap.String("handler", dequeuedJob.HandlerName), zap.Any("error", handlerError))

	if handlerError != nil {
		logger.Log.Error("Error handling job: %v", zap.String("ID", dequeuedJob.ID.String()), zap.String("handler", dequeuedJob.HandlerName), zap.Any("error", handlerError))
	}

	err = q.RemoveProcessed(ctx, dequeuedJob.ID, handlerError)
	if err != nil {
		return fmt.Errorf("error removing processed job: %w", err), waitingMessagePrinted
	}

	return nil, waitingMessagePrinted
}
