// Package app runs the application with needed attributes.
package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gorilla/mux"
	"github.com/katiasuya/audio-conversion-service/internal/auth"
	"github.com/katiasuya/audio-conversion-service/internal/config"
	"github.com/katiasuya/audio-conversion-service/internal/converter"
	"github.com/katiasuya/audio-conversion-service/internal/repository"
	"github.com/katiasuya/audio-conversion-service/internal/server"
	"github.com/katiasuya/audio-conversion-service/internal/storage"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

// Run runs the application service.
func Run() error {
	logger := log.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.DebugLevel)

	var conf config.Config
	conf.Load()
	logger.WithField("package", "app").Infoln("configuration data loaded")

	db, err := repository.NewPostgresDB(&conf)
	if err != nil {
		return fmt.Errorf("can't connect to database: %w", err)
	}
	defer db.Close()

	logger.WithField("package", "app").Infoln("connected to database")
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
	logger.WithField("package", "app").Infoln("cloud storage initialized successfully")

	const maxRequests = 10
	sem := semaphore.NewWeighted(maxRequests)
	converter := converter.New(sem, repo, storage)

	tokenMgr := auth.New(conf.PublicKey, conf.PrivateKey)

	server := server.New(repo, storage, converter, tokenMgr, logger)

	r := mux.NewRouter()
	server.RegisterRoutes(r)

	logger.WithField("package", "app").Infoln("start listening on :8000")

	return http.ListenAndServe(":8000", r)
}
