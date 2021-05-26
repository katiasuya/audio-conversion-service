// Package mycontext provides functions to work with the application context.
package mycontext

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
)

type key int

const (
	userIDKey key = iota
	requestIDKey
)

// Init initializes the logger.
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

// RequestIDFromContext retrieves request id from context.
func RequestIDFromContext(ctx context.Context) (string, bool) {
	requestIDctx, ok := ctx.Value(requestIDKey).(string)
	return requestIDctx, ok
}

// ContextWithUserID adds user id to the context.
func ContextWithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// UserIDFromContext retrieves user id from context.
func UserIDFromContext(ctx context.Context) (string, bool) {
	userIDctx, ok := ctx.Value(userIDKey).(string)
	return userIDctx, ok
}

// LoggerFromContext returns a logger with all possible context.
func LoggerFromContext(ctx context.Context) (newLoggerEntry *log.Entry) {
	newLogger := InitLogger()

	if ctxRqID, ok := ctx.Value(requestIDKey).(string); ok {
		newLoggerEntry = newLogger.WithField("rqID", ctxRqID)
	}
	if ctxUserID, ok := UserIDFromContext(ctx); ok {
		newLoggerEntry = newLoggerEntry.WithField("userID", ctxUserID)
	}

	return newLoggerEntry
}
