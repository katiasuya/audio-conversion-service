package app

import (
	"io/ioutil"
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
	"golang.org/x/sync/semaphore"
)

// Run runs the application service.
func Run() error {
	var conf config.Config
	if err := conf.Load(); err != nil {
		return err
	}

	db, err := repository.NewPostgresDB(&conf)
	if err != nil {
		return err
	}
	defer db.Close()
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

	privateKey, err := getKey(conf.PrivateKeyPath)
	if err != nil {
		return err
	}
	publicKey, err := getKey(conf.PublicKeyPath)
	if err != nil {
		return err
	}
	tokenMgr := auth.New(publicKey, privateKey)

	server := server.New(repo, storage, converter, tokenMgr)

	r := mux.NewRouter()
	server.RegisterRoutes(r)

	return http.ListenAndServe(":8000", r)
}

func getKey(keyPath string) ([]byte, error) {
	file, err := os.Open(keyPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}
