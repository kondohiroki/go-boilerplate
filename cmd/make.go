package cmd

import (
	"os"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/kondohiroki/go-boilerplate/internal/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddGroup(&cobra.Group{ID: "make", Title: "Make:"})
	rootCmd.AddCommand(
		makeMigrationCommand,
	)

	makeMigrationCommand.Flags().StringP("name", "n", "", "(required) Migration name. for example: create_users_table")
	makeMigrationCommand.MarkFlagRequired("name")
}

var makeMigrationCommand = &cobra.Command{
	Use:     "make:migration",
	Short:   "Create a new migration file",
	GroupID: "make",
	Run: func(cmd *cobra.Command, _ []string) {
		// Setup all the required dependencies
		setUpConfig()
		setUpLogger()

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
