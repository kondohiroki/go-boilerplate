package test

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/kondohiroki/go-boilerplate/config"
	"github.com/kondohiroki/go-boilerplate/internal/logger"
	"github.com/kondohiroki/go-boilerplate/internal/router"
	"github.com/valyala/fasthttp"
)

var r *fiber.App

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
	// Set up config
	println("setup config")
	configFile := "../config/config.testing.yaml"
	config.SetConfig(configFile)
	println("setup config done")

	// Set up logger
	println("setup logger")
	logger.InitLogger("zap")
	println("setup logger done")

	// Set up database
	println("setup database")
	migrateTestDatabase()
	println("setup database done")

	println("setup router")
	r = router.NewFiberRouter()
	if r == nil {
		panic("Failed to set up router")
	}
	println("setup router done")
}

// Teardown is called after running tests
func teardown() {

}

func migrateTestDatabase() error {
	return nil
}

// fastHTTPTester returns a new Expect instance to test FastHTTPHandler().
func fastHTTPTester(t *testing.T, handler fasthttp.RequestHandler) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		// Pass requests directly to FastHTTPHandler.
		Client: &http.Client{
			Transport: httpexpect.NewFastBinder(handler),
			Jar:       httpexpect.NewCookieJar(),
		},
		// Report errors using testing.T.
		Reporter: httpexpect.NewAssertReporter(t),
	})
}

func readJSONToString(t *testing.T, filePath string) string {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		t.Errorf("failed to open file: %v", err)
	}
	defer jsonFile.Close()

	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		t.Errorf("failed to read file: %v", err)
	}

	return string(jsonBytes)
}
