package app

import (
	"context"
	"fmt"

	"github.com/katiasuya/audio-conversion-service/internal/config"
	"github.com/katiasuya/audio-conversion-service/internal/converter"
	"github.com/katiasuya/audio-conversion-service/internal/logger"
	"github.com/katiasuya/audio-conversion-service/internal/queue"
	"github.com/katiasuya/audio-conversion-service/internal/repository"
	"github.com/katiasuya/audio-conversion-service/internal/storage"
)

// Run runs the application service.
func RunConverter() error {
	ctx := context.Background()

	conf, err := config.Load()
	if err != nil {
		return fmt.Errorf("can't load configuration: %w", err)
	}
	logger.Info(ctx, "configuration data loaded")

	repo, err := repository.New(&conf.PostgresData)
	if err != nil {
		return err
	}
	defer repo.Close()
	logger.Info(ctx, "connected to database successfully")

	storage, err := storage.NewS3Client(&conf.AWSData)
	if err != nil {
		return fmt.Errorf("can't connect to S3: %w", err)
	}
	logger.Info(ctx, "connected to S3 successfully")

	conn, ch, err := queue.NewRabbitMQClient(&conf.RabbitMQData)
	if err != nil {
		return err
	}
	defer conn.Close()
	defer ch.Close()
	logger.Info(ctx, "connected to RabbitMQ successfully")

	converter := converter.New(repo, storage)
	logger.Info(ctx, "converter initialized successfully")

	queue := queue.New(conf.QueueName, ch, converter)

	return fmt.Errorf("can't process queue messages: %w", queue.ProcessMsgs())
}
