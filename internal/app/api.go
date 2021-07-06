// Package app runs the application with needed attributes.
package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/katiasuya/audio-conversion-service/internal/auth"
	"github.com/katiasuya/audio-conversion-service/internal/config"
	"github.com/katiasuya/audio-conversion-service/internal/logger"
	"github.com/katiasuya/audio-conversion-service/internal/queue"
	"github.com/katiasuya/audio-conversion-service/internal/repository"
	"github.com/katiasuya/audio-conversion-service/internal/server"
	"github.com/katiasuya/audio-conversion-service/internal/storage"
)

// Run runs the application service.
func RunAPI() error {
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

	storage, err := storage.NewS3Client(conf.Bucket, conf.Region, conf.AccessKeyID, conf.SecretAccessKey)
	if err != nil {
		return fmt.Errorf("can't connect to S3: %w", err)
	}
	logger.Info(ctx, "connected to S3 successfully")

	conn, ch, err := queue.NewRabbitMQClient(conf.AmpqURI, conf.QueueName)
	if err != nil {
		return err
	}
	defer conn.Close()
	defer ch.Close()
	logger.Info(ctx, "connnected to RabbitMQ successfully")

	queueMgr := queue.New(conf.QueueName, ch, nil)
	tokenMgr := auth.New(conf.PublicKey, conf.PrivateKey)

	server := server.New(repo, storage, tokenMgr, queueMgr)

	r := mux.NewRouter()
	server.RegisterRoutes(r)

	logger.Info(ctx, "start listening on :8000")
	return http.ListenAndServe(":8000", r)
}
