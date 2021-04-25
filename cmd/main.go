package main

import (
	"github.com/katiasuya/audio-conversion-service/internal/app"
	"github.com/katiasuya/audio-conversion-service/internal/logging"
)

func main() {
	logger := logging.Init()
	logger.WithField("package", "main").Infoln("start listening on :8000")
	logger.WithField("package", "app").Fatalln(app.Run())
}
