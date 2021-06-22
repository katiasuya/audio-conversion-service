// Package logger manages application logging.
package logger

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
)

type key int

const loggerKey key = 0

var defaultLogger = &log.Logger{
	Out:       os.Stdout,
	Formatter: new(log.JSONFormatter),
	Level:     log.InfoLevel,
}

// New returns the default logger.
func New() *log.Logger {
	return defaultLogger
}

// GetFromContext returns logger with all possible context.
func GetFromContext(ctx context.Context) *log.Logger {
	if ctxLogger, ok := ctx.Value(loggerKey).(*log.Logger); ok {
		return ctxLogger
	}
	return defaultLogger
}

// AddToContext adds logger to the context.
func AddToContext(ctx context.Context, logger *log.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// Info logs message at Info level.
func Info(ctx context.Context, msg string) {
	GetFromContext(ctx).Infoln(msg)
}

// Error logs errors at Error level.
func Error(ctx context.Context, err error) {
	GetFromContext(ctx).Errorln(err)
}

// Fatal logs errorss at Fatal level.
func Fatal(ctx context.Context, err error) {
	GetFromContext(ctx).Fatalln(err)
}
