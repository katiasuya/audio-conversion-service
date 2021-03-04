// Package database provides a function to connect a database.
package database

import (
	"database/sql"
	"fmt"

	"github.com/katiasuya/audio-conversion-service/internal/config"
)

//NewPostgresDB creates new database connection.
func NewPostgresDB(c *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.DBName, c.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
