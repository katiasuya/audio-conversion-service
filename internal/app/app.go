package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/katiasuya/audio-conversion-service/internal/config"
	"github.com/katiasuya/audio-conversion-service/internal/repository"
	"github.com/katiasuya/audio-conversion-service/internal/server"
)

// Run runs the application service.
func Run() error {
	var conf config.Config
	err := conf.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := repository.NewPostgresDB(&conf)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.New(db)
	server := server.New(repo)

	r := mux.NewRouter()
	server.RegisterRoutes(r)

	return http.ListenAndServe(":8000", r)
}
