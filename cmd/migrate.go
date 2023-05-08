package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/kondohiroki/go-boilerplate/config"
	"github.com/kondohiroki/go-boilerplate/internal/db/migrations"
	"github.com/kondohiroki/go-boilerplate/internal/db/pgx"
	"github.com/kondohiroki/go-boilerplate/internal/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddGroup(&cobra.Group{ID: "migrate", Title: "Migrate:"})
	rootCmd.AddCommand(
		migrateCommand,
		migrateDownCommand,
		migrateFlushCommand,
		migrateStatusCommand,
	)

	migrateDownCommand.Flags().Int16P("step", "s", 1, "(Optional) step to rollback")
	migrateDownCommand.Example = "  migrate:down"
	migrateDownCommand.Example += "\n  migrate:down -s 2"
}

var migrateCommand = &cobra.Command{
	Use:     "migrate",
	Short:   "Migrate database",
	GroupID: "migrate",
	Run: func(_ *cobra.Command, _ []string) {
		// Setup all the required dependencies
		setUpConfig()
		setUpLogger()
		setUpPostgres()

		if len(migrations.Migrations) == 0 {
			logger.Log.Info("No migrations found")
			os.Exit(0)
		}

		// Initiate context
		ctx := context.Background()

		// Get database connection
		dbConn := pgx.GetPgxPool()

		if dbConn == nil {
			logger.Log.Error("Database connection is nil")
			return
		}

		// Create the migrations table if it doesn't exist
		_, err := dbConn.Exec(
			ctx,
			`CREATE TABLE IF NOT EXISTS migrations (
					id SERIAL PRIMARY KEY,
					migration VARCHAR(255) NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT NOW()
				)`,
		)
		if err != nil {
			logger.Log.Error("Error creating migrations table", zap.Error(err))
			return
		}

		// Get the latest migration that has been applied
		var latestMigration string
		_ = dbConn.QueryRow(
			ctx,
			`SELECT migration FROM migrations ORDER BY id DESC LIMIT 1`,
		).Scan(&latestMigration)

		if latestMigration == "" {
			latestMigration = "0"
		}

		// Check if the latest migration is the last migration
		if latestMigration == migrations.Migrations[len(migrations.Migrations)-1].Name {
			logger.Log.Info("Database is already up-to-date")
			os.Exit(0)
		} else {
			logger.Log.Info("Database is not up-to-date")
			logger.Log.Info("Latest migration: " + latestMigration)
		}

		// Run migrations
		for _, migration := range migrations.Migrations {
			if migration.Name > latestMigration {
				logger.Log.Info("Running migration: " + migration.Name)
				err := migration.Up()
				if err != nil {
					logger.Log.Error("Error running migration", zap.Error(err))
					return
				}

				// Insert the migration into the migrations table
				_, err = dbConn.Exec(
					ctx,
					`INSERT INTO migrations (migration) VALUES ($1)`,
					migration.Name,
				)
				if err != nil {
					logger.Log.Error("Error inserting migration into migrations table", zap.Error(err))
					return
				}
			}
		}
	},
}

var migrateDownCommand = &cobra.Command{
	Use:     "migrate:down",
	Short:   "Rollback database migration",
	GroupID: "migrate",
	Run: func(cmd *cobra.Command, _ []string) {
		// Setup all the required dependencies
		setUpConfig()
		setUpLogger()
		setUpPostgres()

		if len(migrations.Migrations) == 0 {
			logger.Log.Info("No migrations found")
			os.Exit(0)
		}

		// Initiate context
		ctx := context.Background()

		// Get database connection
		dbConn := pgx.GetPgxPool()

		if dbConn == nil {
			logger.Log.Error("Database connection is nil")
			return
		}

		// Get the step flag value default is 1
		step, _ := cmd.Flags().GetInt16("step")

		// Get the applied migrations
		var appliedMigrations []string
		rows, err := dbConn.Query(ctx, `SELECT migration FROM migrations ORDER BY id DESC`)
		if err != nil {
			logger.Log.Error("Error fetching applied migrations", zap.Error(err))
			return
		}
		defer rows.Close()

		for rows.Next() {
			var migration string
			err := rows.Scan(&migration)
			if err != nil {
				logger.Log.Error("Error scanning migration row", zap.Error(err))
				return
			}
			appliedMigrations = append(appliedMigrations, migration)
		}

		if err := rows.Err(); err != nil {
			logger.Log.Error("Error in migration rows", zap.Error(err))
			return
		}

		// Perform the rollback
		for i := 0; i < int(step) && i < len(appliedMigrations); i++ {
			migrationName := appliedMigrations[i]
			logger.Log.Info("Rolling back migration: " + migrationName)

			for _, migration := range migrations.Migrations {
				if migration.Name == migrationName {
					err := migration.Down()
					if err != nil {
						logger.Log.Error("Error rolling back migration", zap.Error(err))
						return
					}

					// Remove the migration from the migrations table
					_, err = dbConn.Exec(
						ctx,
						`DELETE FROM migrations WHERE migration = $1`,
						migration.Name,
					)
					if err != nil {
						logger.Log.Error("Error removing migration from migrations table", zap.Error(err))
						return
					}
					break
				}
			}
		}
	},
}

var migrateFlushCommand = &cobra.Command{
	Use:     "migrate:flush",
	Short:   "Drop all tables in schema",
	GroupID: "migrate",
	Run: func(_ *cobra.Command, _ []string) {
		// Setup all the required dependencies
		setUpConfig()
		setUpLogger()
		setUpPostgres()

		// Initiate context
		ctx := context.Background()

		// Get database connection
		dbConn := pgx.GetPgxPool()

		if dbConn == nil {
			logger.Log.Error("Database connection is nil")
			return
		}

		// Drop all tables in schema
		_, err := dbConn.Exec(
			ctx,
			fmt.Sprintf(
				`DROP SCHEMA IF EXISTS %s CASCADE`,
				config.GetConfig().Postgres.Schema,
			),
		)
		if err != nil {
			logger.Log.Error("Error dropping all tables in schema", zap.Error(err))
			return
		}

		logger.Log.Info("Dropped all tables in schema " + config.GetConfig().Postgres.Schema + " successfully")
	},
}

var migrateStatusCommand = &cobra.Command{
	Use:     "migrate:status",
	Short:   "Display the status of database migrations",
	GroupID: "migrate",
	Run: func(cmd *cobra.Command, _ []string) {
		// Setup all the required dependencies
		setUpConfig()
		setUpLogger()
		setUpPostgres()

		if len(migrations.Migrations) == 0 {
			logger.Log.Info("No migrations found")
			os.Exit(0)
		}

		// Initiate context
		ctx := context.Background()

		// Get database connection
		dbConn := pgx.GetPgxPool()

		if dbConn == nil {
			logger.Log.Error("Database connection is nil")
			return
		}

		// Get the applied migrations
		appliedMigrationsMap := make(map[string]bool)
		rows, err := dbConn.Query(ctx, `SELECT migration FROM migrations`)
		if err != nil {
			logger.Log.Error("Error fetching applied migrations", zap.Error(err))
			return
		}
		defer rows.Close()

		for rows.Next() {
			var migration string
			err := rows.Scan(&migration)
			if err != nil {
				logger.Log.Error("Error scanning migration row", zap.Error(err))
				return
			}
			appliedMigrationsMap[migration] = true
		}

		if err := rows.Err(); err != nil {
			logger.Log.Error("Error in migration rows", zap.Error(err))
			return
		}

		// Display the status of migrations
		fmt.Println("Migration Status:")
		fmt.Println("-----------------")
		for _, migration := range migrations.Migrations {
			status := "Pending"
			if appliedMigrationsMap[migration.Name] {
				status = "Applied"
			}
			fmt.Printf("%s - %s\n", migration.Name, status)
		}
	},
}
