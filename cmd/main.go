package main

import (
	"github.com/katiasuya/audio-conversion-service/internal/app"
	log "github.com/sirupsen/logrus"
)

func main() {
	logger := log.New()
	logger.WithField("package", "main").Infoln("start listening on :8000")
	logger.WithField("package", "app").Fatalln(app.Run())

}
