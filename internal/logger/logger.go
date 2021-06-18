// Package logger manages application logging.
package logger

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
)

type key int

const loggerKey key = 0

func Init() *log.Entry {
	logger := log.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&log.JSONFormatter{})

	return &log.Entry{Logger: logger}
}

// GetFromContext returns logger with all possible context.
func GetFromContext(ctx context.Context) *log.Entry {
	ctxLogger, ok := ctx.Value(loggerKey).(*log.Entry)
	if !ok {
		ctxLogger = Init()
	}

	return ctxLogger
}

// AddToContext adds logger to the context.
func AddToContext(ctx context.Context, logger *log.Entry) context.Context {
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
