package main

import (
	"log"

	"github.com/katiasuya/audio-conversion-service/internal/app"
	"github.com/katiasuya/audio-conversion-service/internal/config"
	_ "github.com/lib/pq"
)

func main() {
	var conf config.Config
	err := conf.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Fatal(app.Run())
}
