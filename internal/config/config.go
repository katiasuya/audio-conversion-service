// Package config represents application configuration.
package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config represents configuration parameters for the application.
type Config struct {
	Host           string `required:"true"`
	Port           int    `required:"true"`
	Username       string `required:"true"`
	Password       string `required:"true"`
	DBName         string `required:"true"`
	SSLMode        string `required:"true"`
	StoragePath    string `required:"true"`
	PrivateKeyPath string `required:"true"`
	PublicKeyPath  string `required:"true"`
}

// Load loads configuration parameters to Config from environment variables.
func (c *Config) Load() error {
	if err := godotenv.Load("../.env"); err != nil {
		return err
	}
	return envconfig.Process("Audio-converter", c)
}
