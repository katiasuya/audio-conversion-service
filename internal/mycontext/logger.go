package mycontext

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
)

const requestIDKey key = 0

// InitLogger initializes the logger.
func InitLogger() *log.Logger {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.DebugLevel)

	return logger
}

// ContextWithRequestID adds request id to the context.
func ContextWithRequestID(ctx context.Context, rqID string) context.Context {
	return context.WithValue(ctx, requestIDKey, rqID)
}

// LoggerFromContext returns alogger with as much context as possible
func LoggerFromContext(ctx context.Context) *log.Logger {
	newLogger := InitLogger()
	if ctx != nil {
		if ctxRqID, ok := ctx.Value(requestIDKey).(string); ok {
			newLogger = newLogger.WithField("rqId", ctxRqID).Logger
		}
	}

	return newLogger
}
