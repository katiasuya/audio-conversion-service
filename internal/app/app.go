package app

import (
	"net/http"

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
	conf.Load()

	db, err := repository.NewPostgresDB(&conf)
	if err != nil {
		return err
	}
	defer db.Close()

	repo := repository.New(db)
	storage := storage.New(conf.StoragePath)

	const maxRequests = 10
	sem := semaphore.NewWeighted(maxRequests)
	converter := converter.New(sem, repo, storage)

	tokenMgr := auth.New([]byte(conf.PublicKeyPath), []byte(conf.PrivateKeyPath))

	server := server.New(repo, storage, converter, tokenMgr)

	r := mux.NewRouter()
	server.RegisterRoutes(r)

	return http.ListenAndServe(":8000", r)
}
