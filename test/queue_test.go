package test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kondohiroki/go-boilerplate/internal/db/model"
	"github.com/kondohiroki/go-boilerplate/internal/db/rdb"
	"github.com/kondohiroki/go-boilerplate/internal/helper/queue"
	"github.com/kondohiroki/go-boilerplate/internal/job"
)

func TestQueue(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue("testing")

	t.Cleanup(func() {
		// Clean up the queue.
		q.Clear(ctx)
	})

	// Test IsEmpty on an empty queue.
	isEmpty, err := q.IsEmpty(ctx)
	require.NoError(t, err, "IsEmpty should not return an error")
	assert.True(t, isEmpty, "IsEmpty should return true on an empty queue")

	// Test Enqueue.
	job1, _ := job.NewJob("ProcessExample", &job.ProcessExample{
		Data: "Sawadeee Kaab!",
	}, 3, 5)

	err = q.Enqueue(ctx, job1)
	require.NoError(t, err, "Enqueue should not return an error")

	// Test Length and IsEmpty after enqueueing a job.
	length, err := q.Length(ctx)
	require.NoError(t, err, "Length should not return an error")
	assert.Equal(t, int64(1), length, "Length should return 1 after enqueueing a job")

	isEmpty, err = q.IsEmpty(ctx)
	require.NoError(t, err)
	assert.False(t, isEmpty)

	// Test Peek.
	peekedJobs, err := q.Peek(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, 1, len(peekedJobs))

	// Test Dequeue.
	dequeuedJob, err := q.Dequeue(ctx, 1*time.Second)
	require.NoError(t, err)
	assert.Equal(t, job1.Payload, dequeuedJob.Payload)

	// Test RemoveProcessed.
	err = q.RemoveProcessed(ctx, dequeuedJob.ID, nil)
	require.NoError(t, err, "expected no error when removing a processed job")

	// Test RetryFailedByJobID on a non-existing job.
	err = q.RetryFailedByJobID(ctx, uuid.New())
	assert.Error(t, err, "expected error when retrying a non-existing job")

	// Test RetryFailedByJobID on an existing job.
	job2, _ := job.NewJob("ProcessExample", &job.ProcessExample{
		Data: "Sawadeee Kaab2!",
	}, 1, 5)

	err = q.Enqueue(ctx, job2)
	require.NoError(t, err)

	dequeuedJob2, err := q.Dequeue(ctx, 1*time.Second)
	require.NoError(t, err)
	assert.Equal(t, job2.Payload, dequeuedJob2.Payload, "job2 payload not equal to dequeued job2 payload")

	// Assume that the job failed.
	q.RemoveProcessed(ctx, dequeuedJob2.ID, fmt.Errorf("assume job is failed and should be retried until max attempts reached and then moved to failed jobs"))

	// Retry the failed job.
	err = q.RetryFailedByJobID(ctx, dequeuedJob2.ID)
	require.NoError(t, err)

	retriedJob2, err := q.Dequeue(ctx, 1*time.Second)
	require.NoError(t, err)
	assert.Equal(t, job2.Payload, retriedJob2.Payload, "job2 payload not equal to retried job2 payload")

	// Test RetryAllFailed.
	job3, _ := job.NewJob("ProcessExample", &job.ProcessExample{
		Data: "Sawadeee Kaab3!",
	}, 3, 5)
	job4, _ := job.NewJob("ProcessExample", &job.ProcessExample{
		Data: "Sawadeee Kaab4!",
	}, 3, 5)
	err = q.Enqueue(ctx, job3, job4)
	require.NoError(t, err)

	dequeuedJob3, err := q.Dequeue(ctx, 1*time.Second)
	require.NoError(t, err)
	assert.Equal(t, job3.Payload, dequeuedJob3.Payload, "job3 payload not equal to dequeued job3 payload")

	dequeuedJob4, err := q.Dequeue(ctx, 1*time.Second)
	require.NoError(t, err)
	assert.Equal(t, job4.Payload, dequeuedJob4.Payload, "job4 payload not equal to dequeued job4 payload")

	_, err = q.RetryAllFailed(ctx)
	require.NoError(t, err, "retry all failed")

	// Test Clear.
	_, err = q.Clear(ctx)
	require.NoError(t, err, "Clear should not return an error")

	length, err = q.Length(ctx)
	require.NoError(t, err, "Length should not return an error")
	assert.Equal(t, int64(0), length, "Length should return 0 after clearing the queue")

	// Test ListQueueKeys.
	job5, _ := job.NewJob("ProcessExample", &job.ProcessExample{
		Data: "Sawadeee Kaab5!",
	}, 3, 5)
	q.Enqueue(ctx, job5)
	keys, err := queue.ListQueueKeys(ctx)
	require.NoError(t, err, "ListQueueKeys should not return an error")
	found := false
	for _, key := range keys {
		key = rdb.AddQueuePrefix(key)
		t.Logf("key: %s", key)
		if key == q.Key {
			found = true
			break
		}
	}
	assert.True(t, found, "queue [testing] key not found in the list of queue keys")

	// Test ListQueueKeysAndLengths.
	queueInfos, err := queue.ListQueueKeysAndLengths(ctx)
	require.NoError(t, err, "ListQueueKeysAndLengths should not return an error")
	found = false
	for _, queueInfo := range queueInfos {
		if queueInfo.Key == q.Key {
			found = true
			assert.Equal(t, int64(1), queueInfo.NumberOfItems, "queue [testing] should have 1 items")
			break
		}
	}
	assert.True(t, found, "queue [testing] key not found in the list of queue keys and lengths")

	err = repo.Job.ResetProcessingJobsToPending(ctx)
	require.NoError(t, err, "ResetProcessingJobsToPending should not return an error")

	_, err = repo.Job.GetJobs(ctx)
	require.NoError(t, err, "GetJobs should not return an error")

	faildJobPayload, _ := sonic.Marshal(map[string]any{
		"test": "test",
	})
	_, err = repo.Job.AddFailedJob(ctx, model.FaildJob{
		JobID:    uuid.New(),
		Queue:    "testing",
		Payload:  faildJobPayload,
		Error:    "test",
		FailedAt: time.Now(),
	})
	require.NoError(t, err, "AddFailedJob should not return an error")

	_, err = repo.Job.GetFailedJobs(ctx)
	require.NoError(t, err, "GetFailedJobs should not return an error")

	_, err = repo.Job.GetUnfinishedJobs(ctx)
	require.NoError(t, err, "GetUnfinishedJobs should not return an error")
}

func TestGetQueues(t *testing.T) {
	type params struct{}

	ctx := context.Background()

	data := job.ProcessExample{
		Data: "Hello World",
	}
	job, _ := job.NewJob("ProcessExample", data, 2, 2)

	testGetQueue := queue.NewQueue("test_get_queue")
	testGetQueue.Enqueue(ctx, job)

	t.Cleanup(func() {
		testGetQueue.Clear(ctx)
	})

	tests := []struct {
		name               string
		params             params
		body               any
		expectedStatusCode int
		expectedSchema     string
		expectedCode       int
		expectedMessage    string
	}{
		{
			name:               "test get all queues",
			body:               "",
			expectedStatusCode: http.StatusOK,
			expectedSchema:     readJSONToString(t, "json_response_schema/get_queues.json"),
			expectedCode:       0,
			expectedMessage:    "OK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := fastHTTPTester(t, r.Handler())

			resp := e.GET("/api/v1/queues").Expect()

			resp.Status(tt.expectedStatusCode)
			resp.JSON().Schema(tt.expectedSchema)
			resp.JSON().Object().Value("response_code").IsEqual(tt.expectedCode)
			resp.JSON().Object().Value("response_message").IsEqual(tt.expectedMessage)

		})
	}
}
