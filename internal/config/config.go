// Package config represents application configuration.
package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config represents configuration parameters for the application.
type Config struct {
	PostgresData
	JWTKeys
	AWSData
	RabbitMQData
}

type PostgresData struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
	SSLMode  string
}

type JWTKeys struct {
	PrivateKey string
	PublicKey  string
}

type AWSData struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Bucket          string
}

type RabbitMQData struct {
	URI       string
	QueueName string
}

// Load loads configuration parameters to Config from environment variables.
func Load() (*Config, error) {
	var conf Config
	err := envconfig.Process("converter", &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
