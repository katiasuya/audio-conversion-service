// Package config represents application configuration.
package config

import (
	"os"
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
	Username string
	Password string
	DBName   string
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
	AmpqURI   string
	QueueName string
}

// Load loads configuration parameters to Config from environment variables.
func (c *Config) Load() {
	c.Host = os.Getenv("POSTGRES_HOST")
	c.Port = os.Getenv("POSTGRES_PORT")
	c.Username = os.Getenv("POSTGRES_USER")
	c.Password = os.Getenv("POSTGRES_PASSWORD")
	c.DBName = os.Getenv("POSTGRES_DB")
	c.SSLMode = os.Getenv("SSLMODE")
	c.PrivateKey = os.Getenv("PRIVATEKEY")
	c.PublicKey = os.Getenv("PUBLICKEY")
	c.AccessKeyID = os.Getenv("AWS_ACCESSKEYID")
	c.SecretAccessKey = os.Getenv("AWS_SECRETACCESSKEY")
	c.Region = os.Getenv("AWS_REGION")
	c.Bucket = os.Getenv("AWS_BUCKET")
	c.AmpqURI = os.Getenv("AMQP_URI")
	c.QueueName = os.Getenv("QUEUE_NAME")
}
