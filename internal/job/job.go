package job

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type JobHandler interface {
	Handle() error
}

// Job represents a job in the queue with a unique ID, queue name, payload, and creation timestamp.
type Job struct {
	ID          uuid.UUID       `json:"id"`
	HandlerName string          `json:"handlerName"`
	Payload     json.RawMessage `json:"payload"`
	CreatedAt   time.Time       `json:"created_at"`
	MaxAttempts int             `json:"max_attempts"`
	Attempts    int             `json:"attempts"`
	Delay       time.Duration   `json:"delay"`
	Errors      []string        `json:"errors"`
}

// NewJob creates a new Job with the given queue name and payload.
func NewJob(handlerName string, payload any, maxAttempts int, delay time.Duration) (*Job, error) {
	jobID := uuid.New()
	createdAt := time.Now()

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Job{
		ID:          jobID,
		HandlerName: handlerName,
		Payload:     payloadBytes,
		MaxAttempts: maxAttempts,
		Attempts:    0,
		Delay:       delay,
		CreatedAt:   createdAt,
	}, nil
}
