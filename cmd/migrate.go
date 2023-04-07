package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/kondohiroki/go-boilerplate/config"
	"github.com/kondohiroki/go-boilerplate/internal/db"
	"github.com/kondohiroki/go-boilerplate/internal/db/migrations"
	"github.com/kondohiroki/go-boilerplate/internal/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(
		dbMigrateCommand,
		dbMigrateFlushCommand,
		dbMigrateCreateCommand,
	)

	dbMigrateCreateCommand.Flags().StringP("name", "n", "", "(required) Migration name. for example: create_users_table")
	dbMigrateCreateCommand.MarkFlagRequired("name")
}

var dbMigrateCommand = &cobra.Command{
	Use:   "db:migrate",
	Short: "Migrate database",
	Run: func(_ *cobra.Command, _ []string) {
		if len(migrations.Migrations) == 0 {
			logger.Log.Info("No migrations found")
			os.Exit(0)
		}

		// Initiate context
		ctx := context.Background()

		// Get database connection
		dbConn := db.GetPgxPool()

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

var dbMigrateFlushCommand = &cobra.Command{
	Use:   "db:migrate:flush",
	Short: "Migrate database and drop all data",
	Run: func(_ *cobra.Command, _ []string) {
		// Initiate context
		ctx := context.Background()

		// Get database connection
		dbConn := db.GetPgxPool()

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

var dbMigrateCreateCommand = &cobra.Command{
	Use:   "db:migrate:create",
	Short: "Create a new migration file",
	Run: func(cmd *cobra.Command, _ []string) {
		migrationFileName, _ := cmd.Flags().GetString("name")

		template := "pkg/template/migration_file.txt"

		// make output with timestamp
		outputFilename := time.Now().Format("20060102150405") + "_" + migrationFileName
		outputPath := "internal/db/migrations/" + outputFilename + ".go"

		// replace migration name
		migrationName := strings.ReplaceAll(migrationFileName, "_", " ")
		migrationName = strcase.ToLowerCamel(migrationName)

		read, err := os.ReadFile(template)
		if err != nil {
			logger.Log.Error("Error reading template file", zap.Error(err))
		}

		newContents := strings.Replace(string(read), "<migration_name>", migrationName, -1)
		newContents = strings.Replace(string(newContents), "<filename>", outputFilename, -1)

		_, err = os.Create(outputPath)
		if err != nil {
			logger.Log.Error("Error creating file", zap.Error(err))
		}

		os.Chmod(outputPath, 0777)

		err = os.WriteFile(outputPath, []byte(newContents), 0)
		if err != nil {
			logger.Log.Error("Error writing to file", zap.Error(err))
		}

	},
}
