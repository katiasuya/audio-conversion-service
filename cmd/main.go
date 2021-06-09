package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/katiasuya/audio-conversion-service/internal/app"
)

func main() {
	logger := log.New()
	logger.WithField("package", "app").Fatalln(app.Run())
}
