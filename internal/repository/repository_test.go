package repository

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/katiasuya/audio-conversion-service/internal/server/model"
	"github.com/tj/assert"
)

var a = &model.AudioInfo{
	ID:       uuid.New().String(),
	Name:     "audio",
	Format:   "mp3",
	Location: uuid.New().String(),
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		// maybe not log
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestGetAudioByID(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}
	defer func() {
		repo.Close()
	}()

	query := `SELECT name, format, location FROM converter.audio WHERE id = (?)`
	rows := sqlmock.NewRows([]string{"name", "format", "location"}).
		AddRow(a.Name, a.Format, a.Location)

	mock.ExpectQuery(query).WithArgs(a.ID).WillReturnRows(rows)

	user, err := repo.GetAudioByID(a.ID)
	assert.NotNil(t, user)
	assert.NoError(t, err)
}
