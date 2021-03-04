package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/katiasuya/audio-conversion-service/internal/app"
	"github.com/kelseyhightower/envconfig"
)

// Config represents configuration parametres for db connection.
type config struct {
	Username string `required:"true"`
	Password string `required:"true"`
	Host     string `required:"true"`
	Port     int    `required:"true"`
}

func main() {
	var conf config
	err := loadConfig(&conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Fatal(app.Run())
}

func loadConfig(c *config) error {
	if err := godotenv.Load("../db.env"); err != nil {
		return err
	}
	return envconfig.Process("Audio-converter", c)
}
