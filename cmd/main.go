package main

import (
	"github.com/katiasuya/audio-conversion-service/internal/app"
	"github.com/katiasuya/audio-conversion-service/internal/mycontext"
)

func main() {
	logger := mycontext.Init()
	logger.WithField("package", "main").Infoln("start listening on :8000")
	logger.WithField("package", "app").Fatalln(app.Run())
}
