package test

import (
	"os"
	"testing"
)

// TestMain is the entry point for running tests
func TestMain(m *testing.M) {
	setup()
	exitCode := m.Run()
	if exitCode == 0 {
		teardown()
	}
	os.Exit(exitCode)
}

// Setup is called before running tests
func setup() {
}

// Teardown is called after running tests
func teardown() {

}

func migrateTestDatabase() error {
	return nil
}
