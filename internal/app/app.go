package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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
	err := conf.Load()
	if err != nil {
		return fmt.Errorf("can't load configuration data: %w", err)
	}
	logger.WithField("package", "app").Infoln("configuration data loaded successfully")

	db, err := repository.NewPostgresDB(&conf)
	if err != nil {
		return fmt.Errorf("can't connect to database: %w", err)
	}
	defer db.Close()
	logger.WithField("package", "app").Infoln("connected to database")

	repo := repository.New(db)
	storage := storage.New(conf.StoragePath)

	const maxRequests = 10
	sem := semaphore.NewWeighted(maxRequests)
	converter := converter.New(sem, repo, storage)

	privateKey, err := getKey(conf.PrivateKeyPath)
	if err != nil {
		return fmt.Errorf("can't load private key: %w", err)
	}
	publicKey, err := getKey(conf.PublicKeyPath)
	if err != nil {
		return fmt.Errorf("can't load public key: %w", err)
	}
	logger.WithField("package", "app").Infoln("keys loaded successfully")
	tokenMgr := auth.New(publicKey, privateKey)

	server := server.New(repo, storage, converter, tokenMgr, logger.WithField("package", "server"))

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
