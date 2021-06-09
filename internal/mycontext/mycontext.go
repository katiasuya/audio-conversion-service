// Package mycontext provides functions to work with the application context.
package mycontext

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type key int

const (
	userIDKey key = iota
	loggerKey
)

// ContextWithUserID adds user id to the context.
func ContextWithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// UserIDFromContext retrieves user id from context.
func UserIDFromContext(ctx context.Context) (string, bool) {
	userIDctx, ok := ctx.Value(userIDKey).(string)
	return userIDctx, ok
}

// ContextWithLogger adds logger to the context.
func ContextWithLogger(ctx context.Context, logger *log.Entry) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// LoggerFromContext returns a logger with all possible context.
func LoggerFromContext(ctx context.Context) (ctxLogger *log.Entry, ok bool) {
	ctxLogger, ok = ctx.Value(loggerKey).(*log.Entry)
	if !ok {
		return
	}

	if ctxUserID, ok := UserIDFromContext(ctx); ok {
		ctxLogger = ctxLogger.WithField("userID", ctxUserID)
	}

	return
}
