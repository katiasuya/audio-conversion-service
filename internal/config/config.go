// Package config represents application configuration.
package config

import (
	"os"
)

// Config represents configuration parameters for the application.
type Config struct {
	Host            string `required:"true"`
	Port            string `required:"true"`
	Username        string `required:"true"`
	Password        string `required:"true"`
	DBName          string `required:"true"`
	SSLMode         string `required:"true"`
	PrivateKey      string `required:"true"`
	PublicKey       string `required:"true"`
	AccessKeyID     string `required:"true"`
	SecretAccessKey string `required:"true"`
	Region          string `required:"true"`
	Bucket          string `required:"true"`
	AmpqUri         string `required:"true"`
	QueueName       string `required:"true"`
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
	c.AmpqUri = os.Getenv("AMQP_URI")
	c.QueueName = os.Getenv("QUEUE_NAME")
}
