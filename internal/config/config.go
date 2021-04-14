// Package config represents application configuration.
package config

import (
	"os"
	"strconv"
)

// Config represents configuration parameters for the application.
type Config struct {
	Host            string `required:"true"`
	Port            int    `required:"true"`
	Username        string `required:"true"`
	Password        string `required:"true"`
	DBName          string `required:"true"`
	SSLMode         string `required:"true"`
	PrivateKeyPath  string `required:"true"`
	PublicKeyPath   string `required:"true"`
	AccessKeyID     string `required:"true"`
	SecretAccessKey string `required:"true"`
	Region          string `required:"true"`
	Bucket          string `required:"true"`
}

// Load loads configuration parameters to Config from environment variables.
func (c *Config) Load() {
	c.Host = os.Getenv("HOST")
	c.Port, _ = strconv.Atoi(os.Getenv("PORT"))
	c.Username = os.Getenv("POSTGRES_USER")
	c.Password = os.Getenv("POSTGRES_PASSWORD")
	c.DBName = os.Getenv("POSTGRES_DB")
	c.SSLMode = os.Getenv("SSLMODE")
	c.StoragePath = os.Getenv("STORAGEPATH")
	c.PrivateKeyPath = os.Getenv("PRIVATEKEYPATH")
	c.PublicKeyPath = os.Getenv("PUBLICKEYPATH")
}
