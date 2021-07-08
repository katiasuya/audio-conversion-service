// Package repository provides the logic for working with database.
package repository

import (
	"database/sql"
	"fmt"

	"github.com/katiasuya/audio-conversion-service/internal/config"
)

// Repository contains a database that the queries will be sent to
// and provides methods to run these queries.
type Repository struct {
	db *sql.DB
}

// New creates a new repository with provided database.
func New(conf *config.PostgresData) (*Repository, error) {
	db, err := NewPostgresClient(conf)
	if err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return &Repository{db}, nil
}

// Close closes db connection.
func (r *Repository) Close() {
	r.db.Close()
}
