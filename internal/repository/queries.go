package repository

import (
	"database/sql"
	"errors"

	"github.com/katiasuya/audio-conversion-service/internal/repository/model"
	"github.com/lib/pq"
)

const codeUniqueViolation = "23505"

//Errors represent database errors.
var (
	ErrNoSuchAudio       = errors.New("the audio with the given id does not exist")
	ErrNoSuchUser        = errors.New("the user with the given username does not exist")
	ErrUserAlreadyExists = errors.New("the user with the given username already exists")
)

// InsertUser inserts the user into table "user".
func (r *Repository) InsertUser(username, password string) (string, error) {
	var userID string
	const insertUserQuery = `INSERT INTO converter."user" (username, password) VALUES ($1, $2) RETURNING id`

	err := r.db.QueryRow(insertUserQuery, username, password).Scan(&userID)
	if err, ok := err.(*pq.Error); ok && err.Code == codeUniqueViolation {
		return "", ErrUserAlreadyExists
	}

	return userID, err
}

// GetIDAndPasswordByUsername retrieves id and hashed password by the given username.
func (r *Repository) GetIDAndPasswordByUsername(username string) (string, string, error) {
	var userID, password string
	const getIDAndPasswordByUsername = `SELECT id, password FROM converter."user" WHERE username=$1;`

	err := r.db.QueryRow(getIDAndPasswordByUsername, username).Scan(&userID, &password)
	if err == sql.ErrNoRows {
		return "", "", ErrNoSuchUser
	}

	return userID, password, err
}

// InsertAudio inserts the audio into table "audio".
func (r *Repository) InsertAudio(name, format, location string) (string, error) {
	var audioID string
	const insertAudio = `INSERT INTO converter.audio (name, format, location) VALUES
	($1, $2, $3) RETURNING id`

	err := r.db.QueryRow(insertAudio, name, format, location).Scan(&audioID)
	return audioID, err
}

// MakeRequest creates the conversion request and returns its id.
func (r *Repository) MakeRequest(name, sourceFormat, targetFormat, location, userID string) (string, error) {
	var requestID string
	const makeConversionRequest = `WITH audio_id AS (INSERT INTO converter.audio (name, format, location) 
	VALUES ($1, $2, $3) RETURNING id)
	INSERT INTO converter.request (user_id, source_id, source_format, target_id, target_format, status)
	SELECT $4, id, $2, NULL, $5, 'queued'
	FROM audio_id RETURNING id;`

	err := r.db.QueryRow(makeConversionRequest, name, sourceFormat, location, userID, targetFormat).Scan(&requestID)
	return requestID, err
}

// UpdateRequest updates the existing conversion request found by its id.
func (r *Repository) UpdateRequest(requestID, status, targetID string) error {
	var dbTargetID sql.NullString
	if targetID != "" {
		dbTargetID = sql.NullString{String: targetID, Valid: true}
	}

	const updateRequest = `UPDATE converter.request SET target_id=$2, status=$3 WHERE id=$1;`
	_, err := r.db.Exec(updateRequest, requestID, dbTargetID, status)

	return err
}

// GetRequestHistory gets the information about user's requests.
func (r *Repository) GetRequestHistory(userID string) ([]model.Request, error) {
	const getUserRequests = `SELECT r.id, a.name, r.source_format, r.target_format, r.created, r.updated, r.status
    FROM converter.request r JOIN converter.audio a ON a.id = r.source_id
    WHERE r.user_id=$1;`

	rows, err := r.db.Query(getUserRequests, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var req model.Request
	var reqs []model.Request
	for rows.Next() {
		err = rows.Scan(&req.ID, &req.AudioName, &req.SourceFormat, &req.TargetFormat, &req.Created, &req.Updated, &req.Status)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}

	return reqs, rows.Err()
}

// GetAudioByID gets information about the audio with the given id.
func (r *Repository) GetAudioByID(audioID string) (model.Audio, error) {
	var id, name, format, location string
	const getAudioByID = `SELECT name, format, location FROM converter.audio WHERE id = $1;`

	err := r.db.QueryRow(getAudioByID, audioID).Scan(&id, &name, &format, &location)
	if err == sql.ErrNoRows {
		return model.Audio{}, ErrNoSuchAudio
	}

	return model.Audio{ID: id, Name: name, Format: format, Location: location}, err
}
