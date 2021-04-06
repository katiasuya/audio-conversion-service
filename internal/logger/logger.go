package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func Init() *log.Logger {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.DebugLevel)

	return logger
}
