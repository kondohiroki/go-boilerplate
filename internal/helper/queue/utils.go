package queue

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/kondohiroki/go-boilerplate/internal/db/rdb"
)

const (
	ERROR_SCANNING_REDIS_KEY    = "error scanning Redis keys: %w"
	ERROR_GETTING_LENGTH_OF_KEY = "error getting length of key %s: %w"
	ERROR_LISTING_QUEUE_KEY     = "error listing queue keys: %w"
	ERROR_DELETING_KEY          = "error deleting key %s: %w"
)

// ListQueueKeys retrieves all queue keys matching the queue key prefix.
func ListQueueKeys(ctx context.Context) ([]string, error) {
	prefix := rdb.GetQueuePrefix()
	rdbClient := rdb.GetRedisClient()
	var keys []string
	var cursor uint64
	var err error

	for {
		var batch []string
		batch, cursor, err = rdbClient.Scan(ctx, cursor, prefix+"*", 10).Result()
		if err != nil {
			return nil, fmt.Errorf(ERROR_SCANNING_REDIS_KEY, err)
		}

		for _, key := range batch {
			if !strings.HasSuffix(key, "_attempt") && !strings.HasSuffix(key, "_failed") {
				key = strings.TrimPrefix(key, prefix+"_")
				keys = append(keys, key)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return keys, nil
}

// ListQueueKeysAndLengths retrieves all queue keys matching the queue key prefix
// and the number of items in each queue.
func ListQueueKeysAndLengths(ctx context.Context) ([]QueueInfo, error) {
	prefix := rdb.GetQueuePrefix()
	rdbClient := rdb.GetRedisClient()
	var keys []string
	var cursor uint64
	var err error

	for {
		var batch []string
		batch, cursor, err = rdbClient.Scan(ctx, cursor, prefix+"*", 50).Result()
		if err != nil {
			return nil, fmt.Errorf(ERROR_SCANNING_REDIS_KEY, err)
		}

		for _, key := range batch {
			if !strings.HasSuffix(key, "_attempt") && !strings.HasSuffix(key, "_failed") {
				keys = append(keys, key)
			}
		}

		if cursor == 0 {
			break
		}
	}

	// Retrieve the length of each queue.
	queueInfos := make([]QueueInfo, 0, len(keys))
	for _, key := range keys {
		length, err := rdbClient.LLen(ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf(ERROR_GETTING_LENGTH_OF_KEY, key, err)
		}

		queueInfo := QueueInfo{
			Key:              key,
			KeyWithoutPrefix: strings.TrimPrefix(key, prefix+"_"),
			NumberOfItems:    length,
		}
		queueInfos = append(queueInfos, queueInfo)
	}

	return queueInfos, nil
}

// List all failed queue keys matching the queue key prefix.
func ListFailedQueueKeys(ctx context.Context) ([]string, error) {
	prefix := rdb.GetQueuePrefix()
	rdbClient := rdb.GetRedisClient()
	var keys []string
	var cursor uint64
	var err error

	for {
		var batch []string
		batch, cursor, err = rdbClient.Scan(ctx, cursor, prefix+"*", 10).Result()
		if err != nil {
			return nil, fmt.Errorf(ERROR_SCANNING_REDIS_KEY, err)
		}

		for _, key := range batch {
			if strings.HasSuffix(key, "_failed") {
				keys = append(keys, key)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return keys, nil
}

// ClearAll deletes all queue keys matching the queue key prefix and returns the total number of items cleared.
func ClearAll(ctx context.Context) (int64, error) {
	queueKeys, err := ListQueueKeys(ctx)
	if err != nil {
		return 0, fmt.Errorf(ERROR_LISTING_QUEUE_KEY, err)
	}

	rdbClient := rdb.GetRedisClient()
	totalCleared := int64(0)
	for _, key := range queueKeys {
		key := rdb.AddQueuePrefix(key)

		// Get the length of the queue before deleting the key.
		length, err := rdbClient.LLen(ctx, key).Result()
		if err != nil {
			return 0, fmt.Errorf(ERROR_GETTING_LENGTH_OF_KEY, key, err)
		}

		// Delete the key.
		err = rdbClient.Del(ctx, key).Err()
		if err != nil {
			return 0, fmt.Errorf(ERROR_DELETING_KEY, key, err)
		}

		// Add the number of items cleared from this key to the total count.
		totalCleared += length
	}

	return totalCleared, nil
}

// FlushAllFailed deletes all failed queue keys matching the queue key prefix and returns the total number of items cleared.
func FlushAllFailed(ctx context.Context) (int64, error) {
	queueKeys, err := ListFailedQueueKeys(ctx)
	if err != nil {
		return 0, fmt.Errorf(ERROR_LISTING_QUEUE_KEY, err)
	}

	rdbClient := rdb.GetRedisClient()
	totalCleared := int64(0)
	for _, key := range queueKeys {
		// Get the length of the queue before deleting the key.
		length, err := rdbClient.LLen(ctx, key).Result()
		if err != nil {
			return 0, fmt.Errorf(ERROR_GETTING_LENGTH_OF_KEY, key, err)
		}

		// Delete the key.
		err = rdbClient.Del(ctx, key).Err()
		if err != nil {
			return 0, fmt.Errorf(ERROR_DELETING_KEY, key, err)
		}

		// Add the number of items cleared from this key to the total count.
		totalCleared += length
	}

	return totalCleared, nil
}

// RemoveJobByID removes a job by ID automatic traversal of all queues.
func RemoveJobOnAnyQueueByID(ctx context.Context, jobID uuid.UUID) (string, error) {
	// Get all the queue keys
	queueKeys, err := ListQueueKeys(ctx)
	if err != nil {
		return "", fmt.Errorf(ERROR_LISTING_QUEUE_KEY, err)
	}

	// Iterate through all queue keys
	for _, key := range queueKeys {
		q := NewQueue(key)
		println(q.Key)

		// Check if the job is in the main queue
		isDeleted, err := q.RemoveJobByID(ctx, jobID)
		if err != nil {
			return "", fmt.Errorf("error removing job with ID %s from queue %s: %w", jobID, q.KeyWithoutPrefix, err)
		}

		if isDeleted {
			return q.KeyWithoutPrefix, nil
		}

	}

	return "", fmt.Errorf("job with ID %s not found in any queue", jobID)
}
