// Package config represents application configuration.
package config

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
}

// Load loads configuration parameters to Config from environment variables.
func (c *Config) Load() {

	//  c.Host = os.Getenv("HOST")
	// 	c.Port = os.Getenv("PORT")
	// 	c.Username = os.Getenv("POSTGRES_USER")
	// 	c.Password = os.Getenv("POSTGRES_PASSWORD")
	// 	c.DBName = os.Getenv("POSTGRES_DB")
	// 	c.SSLMode = os.Getenv("SSLMODE")
	// 	c.PrivateKey = os.Getenv("PRIVATEKEY")
	// 	c.PublicKey = os.Getenv("PUBLICKEY")
	// 	c.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	// 	c.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	// 	c.Region = os.Getenv("AWS_REGION")
	// 	c.Bucket = os.Getenv("AWS_BUCKET")

	c.Host = `0.0.0.0`
	c.Port = `5432`
	c.Username = `postgres`
	c.Password = `postgres`
	c.DBName = `postgres`
	c.SSLMode = ``
	c.AccessKeyID = ``
	c.SecretAccessKey = ``
	c.Region = ``
	c.Bucket = ``
	c.PrivateKey = ``
	c.PublicKey = ``
}
