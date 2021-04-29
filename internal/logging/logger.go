// Package logging is used for handling structured logging in the application.
package logging

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
)

type key int

const loggerKey key = 0

// Init initializes the logger.
func Init() *log.Logger {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.DebugLevel)

	return logger
}

func ContextWithLogger(ctx context.Context, logger *log.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func LoggerFromContext(ctx context.Context) (*log.Logger, bool) {
	logger, ok := ctx.Value(loggerKey).(*log.Logger)
	return logger, ok
}
