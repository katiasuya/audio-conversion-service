package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/katiasuya/audio-conversion-service/internal/logger"
	"github.com/katiasuya/audio-conversion-service/internal/repository/model"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(context.Background(),
			fmt.Errorf("an error '%s' was not expected when opening a stub database connection", err))
	}

	return db, mock
}

func TestGetAudioByID(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}
	defer repo.Close()

	cases := []struct {
		name          string
		expectedAudio model.Audio
		prepare       func(model.Audio)
	}{
		{
			name: "success",
			expectedAudio: model.Audio{
				ID:       "1",
				Name:     "audio",
				Format:   "mp3",
				Location: "location",
			},
			prepare: func(a model.Audio) {
				mock.ExpectQuery("SELECT (.*) FROM converter.audio WHERE id = ?").
					WithArgs(a.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "format", "location"}).
						AddRow(a.ID, a.Name, a.Format, a.Location))
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(tc.expectedAudio)

			gotAudio, err := repo.GetAudioByID(tc.expectedAudio.ID)
			assertNoError(t, fmt.Errorf("an error '%w' was not expected when getting audio by ID", err))

			if tc.expectedAudio != gotAudio {
				t.Errorf("expected user to be %+v, but got %+v", tc.expectedAudio, gotAudio)
			}

			err = mock.ExpectationsWereMet()
			assertNoError(t, fmt.Errorf("there were unfulfilled expectations: %w", err))
		})
	}

}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if errors.Unwrap(err) != nil {
		t.Errorf("No error expected, got %v error", err)
	}
}
