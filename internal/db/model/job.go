package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID          uuid.UUID       `json:"id"`
	Queue       string          `json:"queue"`
	HandlerName string          `json:"handler_name"`
	Payload     json.RawMessage `json:"payload"`
	MaxAttempts int             `json:"max_attempts"`
	Delay       int             `json:"delay"`
	Status      string          `json:"status"` // "pending", "processing", "completed", "failed"
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	FailedJob   []FaildJob      `json:"failed_job"`
}

type FaildJob struct {
	ID       int             `json:"id"`
	JobID    uuid.UUID       `json:"job_id"`
	Queue    string          `json:"queue"`
	Payload  json.RawMessage `json:"payload"`
	Error    string          `json:"error"`
	FailedAt time.Time       `json:"failed_at"`
}
