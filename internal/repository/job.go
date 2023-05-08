package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kondohiroki/go-boilerplate/internal/db/model"
)

type JobRepository interface {
	AddJob(ctx context.Context, job model.Job) (jobID uuid.UUID, err error)
	AddFailedJob(ctx context.Context, job model.FaildJob) (failedJobID int, err error)
	UpdateJobStatus(ctx context.Context, jobID uuid.UUID, status string) error
	ResetProcessingJobsToPending(ctx context.Context) error
	GetJobs(ctx context.Context) ([]model.Job, error)
	GetUnfinishedJobs(ctx context.Context) ([]model.Job, error)
	GetFailedJobs(ctx context.Context) ([]model.FaildJob, error)
	RemoveFailedJob(ctx context.Context, jobID uuid.UUID) error
}

type JobRepositoryImpl struct {
	pgxPool *pgxpool.Pool
}

func NewJobRepository(pgxPool *pgxpool.Pool) JobRepository {
	return &JobRepositoryImpl{
		pgxPool: pgxPool,
	}
}

func (j *JobRepositoryImpl) AddJob(ctx context.Context, job model.Job) (jobID uuid.UUID, err error) {
	tx, err := j.pgxPool.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		INSERT INTO jobs (id, queue, handler_name, payload, max_attempts, delay, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`, job.ID, job.Queue, job.HandlerName, job.Payload, job.MaxAttempts, job.Delay, job.Status, job.CreatedAt, job.UpdatedAt).Scan(&jobID)
	if err != nil {
		return uuid.Nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	return jobID, nil
}

func (j *JobRepositoryImpl) AddFailedJob(ctx context.Context, job model.FaildJob) (failedJobID int, err error) {
	tx, err := j.pgxPool.Begin(ctx)
	if err != nil {
		return failedJobID, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		INSERT INTO failed_jobs (job_id, queue, payload, error, failed_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, job.JobID, job.Queue, job.Payload, job.Error, job.FailedAt).Scan(&failedJobID)
	if err != nil {
		return failedJobID, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return failedJobID, err
	}

	return failedJobID, nil
}

func (j *JobRepositoryImpl) UpdateJobStatus(ctx context.Context, jobID uuid.UUID, status string) error {
	tx, err := j.pgxPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		UPDATE jobs SET status = $1, updated_at = $2 WHERE id = $3
	`, status, time.Now(), jobID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (j *JobRepositoryImpl) ResetProcessingJobsToPending(ctx context.Context) error {
	tx, err := j.pgxPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		UPDATE jobs SET status = 'pending', updated_at = $1 WHERE status = 'processing'
	`, time.Now())
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func handleSelectJob(rows pgx.Rows) ([]model.Job, error) {
	var jobs []model.Job
	for rows.Next() {
		var job model.Job
		err := rows.Scan(&job.ID, &job.Queue, &job.HandlerName, &job.Payload, &job.MaxAttempts, &job.Delay, &job.Status, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	return jobs, rows.Err()
}

func (j *JobRepositoryImpl) GetJobs(ctx context.Context) ([]model.Job, error) {
	rows, err := j.pgxPool.Query(ctx, `
		SELECT id, queue, handler_name, payload, max_attempts, delay, status, created_at, updated_at FROM jobs
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs, err := handleSelectJob(rows)

	return jobs, err
}

func (j *JobRepositoryImpl) GetUnfinishedJobs(ctx context.Context) ([]model.Job, error) {
	rows, err := j.pgxPool.Query(ctx, `
		SELECT id, queue, handler_name, payload, max_attempts, delay, status, created_at, updated_at FROM jobs WHERE status != 'completed' ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs, err := handleSelectJob(rows)

	return jobs, err
}

func (j *JobRepositoryImpl) GetFailedJobs(ctx context.Context) ([]model.FaildJob, error) {
	rows, err := j.pgxPool.Query(ctx, `
		SELECT id, job_id, queue, payload, error, failed_at FROM failed_jobs ORDER BY failed_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []model.FaildJob
	for rows.Next() {
		var job model.FaildJob
		err := rows.Scan(&job.ID, &job.JobID, &job.Queue, &job.Payload, &job.Error, &job.FailedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	return jobs, rows.Err()
}

func (j *JobRepositoryImpl) RemoveFailedJob(ctx context.Context, jobID uuid.UUID) error {
	tx, err := j.pgxPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		DELETE FROM failed_jobs WHERE job_id = $1
	`, jobID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
