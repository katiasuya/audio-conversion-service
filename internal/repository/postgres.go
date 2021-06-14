package repository

import (
	"database/sql"
	"fmt"

	"github.com/katiasuya/audio-conversion-service/internal/config"
)

//NewPostgresDB creates new database connection.
func NewPostgresDB(c *config.Config) (*sql.DB, error) {
	pqInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.DBName, c.SSLMode)

	db, err := sql.Open("postgres", pqInfo)
	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}
