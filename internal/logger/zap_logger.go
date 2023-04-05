package logger

import (
	"os"

	"github.com/kondohiroki/go-boilerplate/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func newZapLogger() *zap.Logger {
	// STEP 1: Get the log level
	zapLogLevel := getZapLogLevel(config.GetConfig().Log.Level)
	stacktraceLogLevel := getZapLogLevel(config.GetConfig().Log.StacktraceLevel)

	// STEP 2: Set up the file writer
	lumberjackLogger := &lumberjack.Logger{
		Filename:   config.GetConfig().Log.FilePath,
		MaxSize:    config.GetConfig().Log.FileSize,     // megabytes
		MaxBackups: config.GetConfig().Log.MaxBackups,   // number of log files
		MaxAge:     config.GetConfig().Log.MaxAge,       // days
		Compress:   config.GetConfig().Log.FileCompress, // disabled by default
	}

	fileWriter := zapcore.AddSync(lumberjackLogger)

	// STEP 3: Set up the encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// STEP 4: Set up the encoder for the file before changing it for the console
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// STEP 5: Change the time format for the console
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// STEP 6: Set up the core
	var core zapcore.Core
	if config.GetConfig().Log.FileEnabled {
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, fileWriter, zapLogLevel),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapLogLevel),
		)
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapLogLevel),
		)
	}

	// STEP 7: Set up the logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(stacktraceLogLevel))

	return logger
}

func getZapLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}
