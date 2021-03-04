package main

import (
	"log"

	"github.com/katiasuya/audio-conversion-service/internal/app"
	"github.com/katiasuya/audio-conversion-service/internal/config"
	"github.com/katiasuya/audio-conversion-service/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	var conf config.Config
	err := conf.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := database.NewPostgresDB(&conf)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	log.Fatal(app.Run())
}
