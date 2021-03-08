// Package repository provides the logic for working with database.
package repository

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

const codeUniqueViolation = "23505"

// Repository represents the database that the queries will be sent to
// and provides methods to communicate with database.
type Repository struct {
	db *sql.DB
}

// New creates a new repository with provided database.
func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// InsertUser inserts the user into users table.
func (r *Repository) InsertUser(username, password string) (string, error) {
	var userID string
	const insertUserQuery = `INSERT INTO converter."user" (username, password) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(insertUserQuery, username, password).Scan(&userID)
	if err, ok := err.(*pq.Error); ok && err.Code == codeUniqueViolation {
		return "", errors.New("the user with the given username already exists")
	}

	return userID, err
}

// GetUserPassword retrieves the database hashed password of a user.
func (r *Repository) GetUserPassword(username string) (string, error) {
	var password string
	const getPasswordByUsername = `SELECT password FROM converter."user" WHERE username=$1;`
	err := r.db.QueryRow(getPasswordByUsername, username).Scan(&password)
	if err == sql.ErrNoRows {
		return "", errors.New("the user with the given username does not exist")
	}

	return password, err
}

// MakeRequest retrieves the database hashed password of a user.
func (r *Repository) MakeRequest(name, sourceFormat, targetFormat, location, userID string) (string, error) {
	var requestID string
	const makeConversionRequest = `WITH audio_id AS (INSERT INTO converter.audio (name, format, location) VALUES
	($1, $2, $3) RETURNING id)
	INSERT INTO converter.request (user_id, source_id, source_format, target_id, target_format, status)
	SELECT $4, id, $2, NULL, $5, 'queued'
	FROM audio_id RETURNING id;`

	err := r.db.QueryRow(makeConversionRequest, name, sourceFormat, location, userID, targetFormat).Scan(&requestID)
	return requestID, err
}
