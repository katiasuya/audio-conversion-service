// Package repository provides the logic for working with database.
package repository

import (
	"database/sql"
	"errors"

	"github.com/katiasuya/audio-conversion-service/internal/server/model"
	"github.com/lib/pq"
)

const codeUniqueViolation = "23505"

//Errors represent database errors.
var (
	ErrNoSuchAudio       = errors.New("the audio with the given id does not exist")
	ErrNoSuchUser        = errors.New("the user with the given username does not exist")
	ErrUserAlreadyExists = errors.New("the user with the given username already exists")
)

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
		return "", ErrUserAlreadyExists
	}

	return userID, err
}

// GetUserPassword retrieves the database hashed password of a user.
func (r *Repository) GetUserPassword(username string) (string, error) {
	var password string
	const getPasswordByUsername = `SELECT password FROM converter."user" WHERE username=$1;`
	err := r.db.QueryRow(getPasswordByUsername, username).Scan(&password)
	if err == sql.ErrNoRows {
		return "", ErrNoSuchUser
	}

	return password, err
}

// MakeRequest creates the conversion request and returns its id.
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

// GetRequestHistory gets the information about user's requests.
func (r *Repository) GetRequestHistory(userID string) ([]model.RequestInfo, error) {
	const getUserRequests = `SELECT r.id, a.name, source_format, target_format, r.created, r.updated, r.status
    FROM converter.request r JOIN converter.audio a ON a.id = r.source_id
    WHERE r.user_id=$1;`

	rows, err := r.db.Query(getUserRequests, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var req model.RequestInfo
	var reqs []model.RequestInfo
	for rows.Next() {
		err = rows.Scan(&req.ID, &req.AudioName, &req.SourceFormat, &req.TargetFormat, &req.Created, &req.Updated, &req.Status)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}

	return reqs, rows.Err()
}

// GetAudioByID gets the information about the audio with the given id.
func (r *Repository) GetAudioByID(id string) (model.AudioInfo, error) {
	var name, format, location string
	const getAudioByID = `SELECT a.name, a.format, a.location FROM converter.audio  a WHERE id=$1;`

	err := r.db.QueryRow(getAudioByID, id).Scan(&name, &format, &location)
	if err == sql.ErrNoRows {
		return model.AudioInfo{}, ErrNoSuchAudio
	}

	return model.AudioInfo{Name: name, Format: format, Location: location}, err
}
