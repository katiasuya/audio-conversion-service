package main

import (
	"context"
	"fmt"

	"github.com/katiasuya/audio-conversion-service/internal/app"
	"github.com/katiasuya/audio-conversion-service/internal/logger"
)

func main() {
	err := app.RunAPI()
	logger.Fatal(context.Background(), fmt.Errorf("API failed to start: %w", err))
}
