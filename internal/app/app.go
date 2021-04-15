package app

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gorilla/mux"
	"github.com/katiasuya/audio-conversion-service/internal/auth"
	"github.com/katiasuya/audio-conversion-service/internal/config"
	"github.com/katiasuya/audio-conversion-service/internal/converter"
	"github.com/katiasuya/audio-conversion-service/internal/logger"
	"github.com/katiasuya/audio-conversion-service/internal/repository"
	"github.com/katiasuya/audio-conversion-service/internal/server"
	"github.com/katiasuya/audio-conversion-service/internal/storage"
	"golang.org/x/sync/semaphore"
)

// Run runs the application service.
func Run() error {
	logger := logger.Init()

	var conf config.Config
	conf.Load()
  logger.WithField("package", "app").Infoln("configuration data loaded successfully")

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
	uploader := s3manager.NewUploader(sess)
	svc := s3.New(sess)
	storage := storage.New(svc, conf.Bucket, uploader)

	const maxRequests = 10
	sem := semaphore.NewWeighted(maxRequests)
	converter := converter.New(sem, repo, storage)

	tokenMgr := auth.New(conf.PublicKey, conf.PrivateKey)
  logger.WithField("package", "app").Infoln("keys loaded successfully")

	server := server.New(repo, storage, converter, tokenMgr, logger.WithField("package", "server"))

	r := mux.NewRouter()
	server.RegisterRoutes(r)

	return http.ListenAndServe(":8000", r)
}
