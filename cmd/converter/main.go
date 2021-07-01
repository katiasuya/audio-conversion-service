package main

import (
	"context"
	"fmt"

	"github.com/katiasuya/audio-conversion-service/internal/app"
	"github.com/katiasuya/audio-conversion-service/internal/logger"
)

func main() {
	logger.Fatal(context.Background(), fmt.Errorf("converter failed to start: %w", app.RunConverter()))
}
