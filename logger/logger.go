package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func init() {
	// Configure Zap logger
	fmt.Println("I am from Logger init")
	loggerCfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel), // Set log level (Info)
		OutputPaths:      []string{"stdout"},                      // Log to stdout
		ErrorOutputPaths: []string{"stderr"},                      // Log errors to stderr
	}

	// Override the log level based on the environment (release or debug)
	if os.Getenv("GIN_MODE") == "debug" {
		loggerCfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	var err error
	Logger, err = loggerCfg.Build()
	if err != nil {
		panic("Failed to initialize Zap logger: " + err.Error())
	}

	// Replace the standard logger with Zap logger
	zap.ReplaceGlobals(Logger)
}
