package repository

import (
	"database/sql"
	"fmt"

	"github.com/katiasuya/audio-conversion-service/internal/config"
)

// NewPostgresClient creates new postgres connection.
func NewPostgresClient(c *config.PostgresData) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DB, c.SSLMode)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}
