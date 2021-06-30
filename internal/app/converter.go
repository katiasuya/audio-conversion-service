package app

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

	var conf config.Config
	conf.Load()
	logger.Info(ctx, "configuration data loaded")

	db, err := repository.NewPostgresClient(&conf)
	if err != nil {
		return fmt.Errorf("can't connect to database: %w", err)
	}
	defer db.Close()
	logger.Info(ctx, "connected to database")

	repo := repository.New(db)

	sess, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(conf.Region),
			Credentials: credentials.NewStaticCredentials(conf.AccessKeyID, conf.SecretAccessKey, ""),
		},
	)
	if err != nil {
		return fmt.Errorf("can't create new session: %w", err)
	}

	uploader := s3manager.NewUploader(sess)
	svc := s3.New(sess)
	storage := storage.New(svc, conf.Bucket, uploader)
	logger.Info(ctx, "cloud storage initialized successfully")

	conn, ch, err := queue.NewRabbitMQClient(conf.AmpqURI, conf.QueueName)
	if err != nil {
		return err
	}
	defer conn.Close()
	defer ch.Close()
	logger.Info(ctx, "connnected to RabbitMQ successfully")

	converter := converter.New(repo, storage)
	logger.Info(ctx, "converter initialized successfully")

	queue := queue.New(conf.QueueName, ch, converter)

	return fmt.Errorf("can't process queue messages: %w", queue.ProcessMsgs())
}
