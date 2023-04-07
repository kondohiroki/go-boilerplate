package logger_test

import (
	"testing"

	"github.com/kondohiroki/go-boilerplate/config"
	. "github.com/kondohiroki/go-boilerplate/internal/logger"
	"go.uber.org/zap/zapcore"
)

func init() {
	configFile := "../../config/config.testing.yaml"
	config.SetConfig(configFile)
}

func TestInitLogger(t *testing.T) {
	tests := []struct {
		name          string
		logDriver     string
		logLevel      string
		fileEnabled   bool
		expectedLevel zapcore.Level
	}{
		{
			name:          "Initialize logger with info level and no file",
			logDriver:     "zap",
			logLevel:      "info",
			fileEnabled:   false,
			expectedLevel: zapcore.InfoLevel,
		},
		{
			name:          "Initialize logger with debug level and no file",
			logDriver:     "zap",
			logLevel:      "debug",
			fileEnabled:   false,
			expectedLevel: zapcore.DebugLevel,
		},
		{
			name:          "Initialize logger with warn level and no file",
			logDriver:     "zap",
			logLevel:      "warn",
			fileEnabled:   false,
			expectedLevel: zapcore.WarnLevel,
		},
		{
			name:          "Initialize logger with error level and no file",
			logDriver:     "zap",
			logLevel:      "error",
			fileEnabled:   false,
			expectedLevel: zapcore.ErrorLevel,
		},
		{
			name:          "Initialize logger with fatal level and no file",
			logDriver:     "zap",
			logLevel:      "fatal",
			fileEnabled:   false,
			expectedLevel: zapcore.FatalLevel,
		},
		{
			name:          "Initialize logger with panic level and no file",
			logDriver:     "zap",
			logLevel:      "panic",
			fileEnabled:   false,
			expectedLevel: zapcore.PanicLevel,
		},
		{
			name:          "Initialize logger with info level and file",
			logDriver:     "zap",
			logLevel:      "info",
			fileEnabled:   true,
			expectedLevel: zapcore.InfoLevel,
		},
		{
			name:          "Initialize logger with debug level and file",
			logDriver:     "zap",
			logLevel:      "debug",
			fileEnabled:   true,
			expectedLevel: zapcore.DebugLevel,
		},
		{
			name:          "Initialize logger with warn level and file",
			logDriver:     "zap",
			logLevel:      "warn",
			fileEnabled:   true,
			expectedLevel: zapcore.WarnLevel,
		},
		{
			name:          "Initialize logger with error level and file",
			logDriver:     "zap",
			logLevel:      "error",
			fileEnabled:   true,
			expectedLevel: zapcore.ErrorLevel,
		},
		{
			name:          "Initialize logger with fatal level and file",
			logDriver:     "zap",
			logLevel:      "fatal",
			fileEnabled:   true,
			expectedLevel: zapcore.FatalLevel,
		},
		{
			name:          "Initialize logger with panic level and file",
			logDriver:     "zap",
			logLevel:      "panic",
			fileEnabled:   true,
			expectedLevel: zapcore.PanicLevel,
		},
		{
			name:          "Initialize logger with invalid level",
			logDriver:     "zap",
			logLevel:      "invalid_level",
			fileEnabled:   false,
			expectedLevel: zapcore.InfoLevel, // Default level expected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Update config
			cfg := config.GetConfig()
			cfg.Log.Level = tt.logLevel
			cfg.Log.FileEnabled = tt.fileEnabled

			InitLogger(tt.logDriver)

			if Log == nil {
				t.Fatal("Expected logger to be initialized, but it is nil")
			}

			logLevel := Log.Core().Enabled(tt.expectedLevel)
			if !logLevel {
				t.Errorf("Expected log level to be %v, but it's not enabled", tt.expectedLevel)
			}
		})
	}
}
