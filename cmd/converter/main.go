package main

import (
	"context"

	"github.com/katiasuya/audio-conversion-service/internal/app"
	"github.com/katiasuya/audio-conversion-service/internal/logger"
)

func main() {
	logger.Fatal(context.Background(), app.RunConverter())
}
