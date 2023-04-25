package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kondohiroki/go-boilerplate/internal/db/rdb"
	"github.com/kondohiroki/go-boilerplate/internal/helper/queue"
	"github.com/kondohiroki/go-boilerplate/internal/job"
)

func TestQueue(t *testing.T) {
	q := queue.NewQueue("testing")

	t.Cleanup(func() {
		// Clean up the queue.
		q.Clear()
	})

	// Test IsEmpty on an empty queue.
	isEmpty, err := q.IsEmpty()
	require.NoError(t, err, "IsEmpty should not return an error")
	assert.True(t, isEmpty, "IsEmpty should return true on an empty queue")

	// Test Enqueue.
	job1, _ := job.NewJob("ProcessExample", &job.ProcessExample{
		Data: "Sawadeee Kaab!",
	}, 3, time.Second*5)

	err = q.Enqueue(job1)
	require.NoError(t, err, "Enqueue should not return an error")

	// Test Length and IsEmpty after enqueueing a job.
	length, err := q.Length()
	require.NoError(t, err, "Length should not return an error")
	assert.Equal(t, int64(1), length, "Length should return 1 after enqueueing a job")

	isEmpty, err = q.IsEmpty()
	require.NoError(t, err)
	assert.False(t, isEmpty)

	// Test Peek.
	peekedJobs, err := q.Peek(1)
	require.NoError(t, err)
	assert.Equal(t, 1, len(peekedJobs))

	// Test Dequeue.
	dequeuedJob, err := q.Dequeue(1 * time.Second)
	require.NoError(t, err)
	assert.Equal(t, job1.Payload, dequeuedJob.Payload)

	// Test RemoveProcessed.
	err = q.RemoveProcessed(dequeuedJob.ID, nil)
	require.NoError(t, err, "expected no error when removing a processed job")

	// Test RetryFailedByJobID on a non-existing job.
	err = q.RetryFailedByJobID(uuid.New())
	assert.Error(t, err, "expected error when retrying a non-existing job")

	// Test RetryFailedByJobID on an existing job.
	job2, _ := job.NewJob("ProcessExample", &job.ProcessExample{
		Data: "Sawadeee Kaab2!",
	}, 1, time.Second*5)

	err = q.Enqueue(job2)
	require.NoError(t, err)

	dequeuedJob2, err := q.Dequeue(1 * time.Second)
	require.NoError(t, err)
	assert.Equal(t, job2.Payload, dequeuedJob2.Payload, "job2 payload not equal to dequeued job2 payload")

	// Assume that the job failed.
	q.RemoveProcessed(dequeuedJob2.ID, fmt.Errorf("assume job is failed and should be retried until max attempts reached and then moved to failed jobs"))

	// Retry the failed job.
	err = q.RetryFailedByJobID(dequeuedJob2.ID)
	require.NoError(t, err)

	retriedJob2, err := q.Dequeue(1 * time.Second)
	require.NoError(t, err)
	assert.Equal(t, job2.Payload, retriedJob2.Payload, "job2 payload not equal to retried job2 payload")

	// Test RetryAllFailed.
	job3, _ := job.NewJob("ProcessExample", &job.ProcessExample{
		Data: "Sawadeee Kaab3!",
	}, 3, time.Second*5)
	job4, _ := job.NewJob("ProcessExample", &job.ProcessExample{
		Data: "Sawadeee Kaab4!",
	}, 3, time.Second*5)
	err = q.Enqueue(job3, job4)
	require.NoError(t, err)

	dequeuedJob3, err := q.Dequeue(1 * time.Second)
	require.NoError(t, err)
	assert.Equal(t, job3.Payload, dequeuedJob3.Payload, "job3 payload not equal to dequeued job3 payload")

	dequeuedJob4, err := q.Dequeue(1 * time.Second)
	require.NoError(t, err)
	assert.Equal(t, job4.Payload, dequeuedJob4.Payload, "job4 payload not equal to dequeued job4 payload")

	_, err = q.RetryAllFailed()
	require.NoError(t, err, "retry all failed")

	// Test Clear.
	_, err = q.Clear()
	require.NoError(t, err, "Clear should not return an error")

	length, err = q.Length()
	require.NoError(t, err, "Length should not return an error")
	assert.Equal(t, int64(0), length, "Length should return 0 after clearing the queue")

	// Test ListQueueKeys.
	job5, _ := job.NewJob("ProcessExample", &job.ProcessExample{
		Data: "Sawadeee Kaab5!",
	}, 3, time.Second*5)
	q.Enqueue(job5)
	keys, err := queue.ListQueueKeys()
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
	queueInfos, err := queue.ListQueueKeysAndLengths()
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
}

func TestGetQueues(t *testing.T) {
	type params struct{}

	data := job.ProcessExample{
		Data: "Hello World",
	}
	job, _ := job.NewJob("ProcessExample", data, 2, 2*time.Second)

	testGetQueue := queue.NewQueue("test_get_queue")
	testGetQueue.Enqueue(job)

	t.Cleanup(func() {
		testGetQueue.Clear()
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
