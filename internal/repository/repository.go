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
func (r *Repository) InsertUser(username, password string) (err error) {
	const insertUserQuery = `INSERT INTO converter."user" (username, password) VALUES ($1, $2)`
	_, err = r.db.Exec(insertUserQuery, username, password)
	if err, ok := err.(*pq.Error); ok && err.Code == codeUniqueViolation {
		return errors.New("the user with the given username already exists")
	}

	return err
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
