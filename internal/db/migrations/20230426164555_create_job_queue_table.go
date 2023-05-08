package migrations

import (
	"context"

	"github.com/kondohiroki/go-boilerplate/internal/db/pgx"
)

func init() {
	Migrations = append(Migrations, createJobQueueTable)
}

var createJobQueueTable = &Migration{
	Name: "20230426164555_create_job_queue_table",
	Up: func() error {
		_, err := pgx.GetPgxPool().Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS jobs (
			"id" UUID PRIMARY KEY,
			"queue" VARCHAR(255),
			"handler_name" VARCHAR(255),
			"payload" JSONB,
			"max_attempts" INTEGER DEFAULT 1,
			"delay" int,
			"status" VARCHAR(255),
			"created_at" TIMESTAMPTZ DEFAULT NOW(),
			"updated_at" TIMESTAMPTZ DEFAULT NOW()
		  );

		  COMMENT ON COLUMN jobs.status IS 'The status of the job, which can be one of the following: pending, processing, completed, or failed.';

		  CREATE INDEX IF NOT EXISTS idx_jobs_queue ON jobs (queue);
		  CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs (status);
		  
		  CREATE TABLE IF NOT EXISTS failed_jobs (
			"id" SERIAL PRIMARY KEY,
			"job_id" UUID UNIQUE,
			"queue" VARCHAR(255),
			"payload" JSONB,
			"error" TEXT,
			"failed_at" TIMESTAMPTZ
		  );

		  CREATE INDEX IF NOT EXISTS idx_failed_jobs_queue ON failed_jobs (queue);
		`)

		if err != nil {
			return err
		}
		return nil

	},
	Down: func() error {
		_, err := pgx.GetPgxPool().Exec(context.Background(), `
			DROP TABLE IF EXISTS failed_jobs;
			DROP TABLE IF EXISTS jobs;
		`)
		if err != nil {
			return err
		}

		return nil
	},
}
