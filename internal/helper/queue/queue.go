package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"github.com/kondohiroki/go-boilerplate/internal/db/rdb"
	"github.com/kondohiroki/go-boilerplate/internal/job"
	"github.com/kondohiroki/go-boilerplate/internal/logger"
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
	}
}

// Adds an item to the source list (the end of the queue).
func (q *Queue) Enqueue(jobs ...*job.Job) error {
	rdbClient := rdb.GetRedisClient()

	for _, job := range jobs {
		jobBytes, err := sonic.Marshal(job)
		if err != nil {
			return err
		}

		err = rdbClient.LPush(context.Background(), q.Key, jobBytes).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

// Removes an item from the source list (the start of the queue) and adds it to the destkey list (temporary storage location).
func (q *Queue) Dequeue(timeout time.Duration) (*job.Job, error) {
	sourceKey := q.Key
	destKey := q.Key + "_tmp"
	rdbClient := rdb.GetRedisClient()

	result, err := rdbClient.BLMove(context.Background(), sourceKey, destKey, "RIGHT", "LEFT", timeout).Result()

	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var job job.Job
	err = sonic.Unmarshal([]byte(result), &job)
	if err != nil {
		return nil, err
	}

	// Check if the job has reached the maximum number of attempts
	if job.MaxAttempts > 0 && job.Attempts >= job.MaxAttempts {
		// Remove the job from the temporary list, it has reached the maximum number of attempts
		_, _ = rdbClient.LRem(context.Background(), destKey, 1, result).Result()
		return nil, fmt.Errorf("job %s reached the maximum number of attempts (%d)", job.ID, job.MaxAttempts)
	}

	// Increment the number of attempts and update the job in the temporary list
	job.Attempts++
	jobBytes, err := sonic.Marshal(&job)
	if err != nil {
		return nil, err
	}

	// Remove the job from the temporary list
	_, err = rdbClient.LRem(context.Background(), destKey, 1, result).Result()
	if err != nil {
		return nil, err
	}

	// Add the updated job to the temporary list
	err = rdbClient.LPush(context.Background(), destKey, jobBytes).Err()
	if err != nil {
		return nil, err
	}

	return &job, nil
}

// Removes the processed item with the given job ID from the destkey list (temporary storage location).
func (q *Queue) RemoveProcessed(jobID uuid.UUID, jobError error) error {
	destKey := q.Key + "_tmp"
	sourceKey := q.Key
	failedJobsKey := q.Key + "_failed"
	rdbClient := rdb.GetRedisClient()

	queues := rdbClient.LRange(context.Background(), destKey, 0, -1)
	for _, queueItem := range queues.Val() {
		var job job.Job
		if err := sonic.Unmarshal([]byte(queueItem), &job); err != nil {
			return err
		}

		if job.ID == jobID {
			if err := removeJobFromList(rdbClient, destKey, queueItem); err != nil {
				return err
			}

			if jobError != nil {
				job.Errors = append(job.Errors, jobError.Error())
				return handleFailedJob(rdbClient, job, sourceKey, failedJobsKey)
			}
			break
		}
	}

	return nil
}

// Remove the job from the temporary list if it was processed successfully
func removeJobFromList(rdbClient *redis.Client, destKey string, queueItem string) error {
	if _, err := rdbClient.LRem(context.Background(), destKey, 1, queueItem).Result(); err != nil {
		return err
	}
	return nil
}

// If the job failed, add it to the failed_jobs list
func handleFailedJob(rdbClient *redis.Client, job job.Job, sourceKey string, failedJobsKey string) error {
	if job.MaxAttempts == 0 || job.Attempts < job.MaxAttempts {
		time.Sleep(job.Delay)
		jobBytes, err := sonic.Marshal(&job)
		if err != nil {
			return err
		}
		err = rdbClient.LPush(context.Background(), sourceKey, jobBytes).Err()
		if err != nil {
			return err
		}
	} else {
		if err := addJobToFailedList(rdbClient, job, failedJobsKey); err != nil {
			return err
		}
	}
	return nil
}

func addJobToFailedList(rdbClient *redis.Client, job job.Job, failedJobsKey string) error {
	logger.Log.Info("Job has reached the maximum number of attempts. It will be added to the failed_jobs list", zap.String("job_id", job.ID.String()))
	jobBytes, err := sonic.Marshal(&job)
	if err != nil {
		return err
	}
	err = rdbClient.LPush(context.Background(), failedJobsKey, jobBytes).Err()
	if err != nil {
		return err
	}
	return nil
}

// RetryFailedByJobID moves a failed item with the given job ID from the destkey list back to the source list for retrying.
func (q *Queue) RetryFailedByJobID(jobID uuid.UUID) error {
	failedJobsKey := q.Key + "_failed"
	destkey := q.Key
	rdbClient := rdb.GetRedisClient()

	queues, err := rdbClient.LRange(context.Background(), failedJobsKey, 0, -1).Result()
	if err != nil {
		return err
	}

	found, err := findAndRetryJob(rdbClient, queues, jobID, failedJobsKey, destkey)
	if err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("job with ID %s not found in temporary list", jobID)
	}

	return nil
}

func findAndRetryJob(rdbClient *redis.Client, queues []string, jobID uuid.UUID, failedJobsKey, destkey string) (bool, error) {
	for _, failedItem := range queues {
		var job job.Job
		if err := sonic.Unmarshal([]byte(failedItem), &job); err != nil {
			return false, err
		}

		if job.ID == jobID {
			if err := resetAndMoveJob(rdbClient, job, failedJobsKey, destkey, failedItem); err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}

func resetAndMoveJob(rdbClient *redis.Client, job job.Job, failedJobsKey, destkey, failedItem string) error {
	job.Attempts = 0
	updatedItem, err := sonic.Marshal(&job)
	if err != nil {
		return err
	}

	if _, err := rdbClient.LRem(context.Background(), failedJobsKey, 1, failedItem).Result(); err != nil {
		return err
	}

	err = rdbClient.RPush(context.Background(), destkey, updatedItem).Err()
	if err != nil {
		_ = rdbClient.LPush(context.Background(), failedJobsKey, failedItem).Err()
		return err
	}

	return nil
}

// Moves all failed items from the destkey list back to the source list for retrying.
func (q *Queue) RetryAllFailed() (int, error) {
	failedJobsKey := q.Key + "_failed"
	destkey := q.Key
	rdbClient := rdb.GetRedisClient()

	count := 0

	for {
		// Remove the failed item from the temporary list.
		failedItem, err := rdbClient.LPop(context.Background(), failedJobsKey).Result()
		if err == redis.Nil {
			// No more items left in the temporary list, break the loop.
			break
		}
		if err != nil {
			return count, err
		}

		// Reset the attempts counter
		var job job.Job
		if err := sonic.Unmarshal([]byte(failedItem), &job); err != nil {
			return count, err
		}
		job.Attempts = 0
		updatedItem, err := sonic.Marshal(&job)
		if err != nil {
			return count, err
		}

		// Add the failed item back to the source list with the reset attempts counter.
		err = rdbClient.RPush(context.Background(), destkey, updatedItem).Err()
		if err != nil {

			// Add the failed item back to the temporary list in case of an RPush error.
			_ = rdbClient.LPush(context.Background(), failedJobsKey, failedItem).Err()
			return count, err
		}

		count++
	}

	return count, nil
}

// Returns the current length of the source list (the number of items in the queue).
func (q *Queue) Length() (int64, error) {
	rdbClient := rdb.GetRedisClient()

	length, err := rdbClient.LLen(context.Background(), q.Key).Result()
	if err != nil {
		return 0, err
	}

	return length, nil
}

// IsEmpty checks if the source list (queue) is empty.
func (q *Queue) IsEmpty() (bool, error) {
	rdbClient := rdb.GetRedisClient()

	length, err := rdbClient.LLen(context.Background(), q.Key).Result()
	if err != nil {
		return false, err
	}

	return length == 0, nil
}

// Clear removes all items from the source list (queue).
func (q *Queue) Clear() (int64, error) {
	rdbClient := rdb.GetRedisClient()

	// Get the length of the queue before deleting the key.
	length, err := rdbClient.LLen(context.Background(), q.Key).Result()
	if err != nil {
		return 0, fmt.Errorf("error getting length of key %s: %w", q.Key, err)
	}

	// Remove all items from the source list.
	_, err = rdbClient.Del(context.Background(), q.Key).Result()
	if err != nil {
		return 0, err
	}

	return length, nil
}

// RemoveJobByID removes the job with the matching job ID from the source list.
func (q *Queue) RemoveJobByID(jobID uuid.UUID) (bool, error) {
	rdbClient := rdb.GetRedisClient()

	queues, err := rdbClient.LRange(context.Background(), q.Key, 0, -1).Result()
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
			if _, err := rdbClient.LRem(context.Background(), q.Key, 1, queueItem).Result(); err != nil {
				return false, err
			}

			return true, nil
		}
	}

	return false, nil
}

// RemoveFailedByID removes the failed item with the matching job ID from the failed list.
func (q *Queue) RemoveFailedByID(jobID uuid.UUID) error {
	failedJobsKey := q.Key + "_failed"
	rdbClient := rdb.GetRedisClient()

	queues, err := rdbClient.LRange(context.Background(), failedJobsKey, 0, -1).Result()
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
			if _, err := rdbClient.LRem(context.Background(), failedJobsKey, 1, failedItem).Result(); err != nil {
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
func (q *Queue) RemoveAllFailed() (int64, error) {
	rdbClient := rdb.GetRedisClient()

	// Get the length of the queue before deleting the key.
	length, err := rdbClient.LLen(context.Background(), q.Key).Result()
	if err != nil {
		return 0, fmt.Errorf("error getting length of key %s: %w", q.Key, err)
	}

	// Remove all items from the failed list.
	_, err = rdbClient.Del(context.Background(), q.Key+"_failed").Result()
	if err != nil {
		return 0, err
	}

	return length, nil
}

// Peek returns the first N items in the source list without removing them.
func (q *Queue) Peek(count int64) ([]interface{}, error) {
	rdbClient := rdb.GetRedisClient()

	rawItems, err := rdbClient.LRange(context.Background(), q.Key, 0, count-1).Result()
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
			_, printed := processJob(q, handlerMap, waitingMessagePrinted)
			waitingMessagePrinted = printed
		}
	}
}

func processJob(q *Queue, handlerMap job.HandlerMap, waitingMessagePrinted bool) (error, bool) {
	dequeuedJob, err := q.Dequeue(time.Second * 5)
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
		logger.Log.Error("Handler not found", zap.String("ID", dequeuedJob.ID.String()), zap.String("handler", dequeuedJob.HandlerName))
		err := fmt.Errorf("handler not found: %v", dequeuedJob.HandlerName)
		q.RemoveProcessed(dequeuedJob.ID, err)
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

	if err != nil {
		logger.Log.Error("Error handling job: %v", zap.String("ID", dequeuedJob.ID.String()), zap.String("handler", dequeuedJob.HandlerName), zap.Error(err))
	}

	err = q.RemoveProcessed(dequeuedJob.ID, handlerError)
	if err != nil {
		logger.Log.Info("Error removing processed job: %v", zap.String("ID", dequeuedJob.ID.String()), zap.String("handler", dequeuedJob.HandlerName), zap.Error(err))
	}

	return nil, waitingMessagePrinted
}
