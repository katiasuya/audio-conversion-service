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
			assertNoError(t, fmt.Errorf("error when getting audio by ID: '%w'", err))

			if tc.expectedAudio != gotAudio {
				t.Errorf("expected user to be %+v, but got %+v", tc.expectedAudio, gotAudio)
			}

			err = mock.ExpectationsWereMet()
			assertNoError(t, fmt.Errorf("there were unfulfilled expectations: %w", err))
		})
	}
}
func TestInsertUser(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}
	defer repo.Close()

	cases := []struct {
		name         string
		expectedUser model.User
		prepare      func(model.User)
	}{
		{
			name: "success",
			expectedUser: model.User{
				ID:       "1",
				Username: "username",
				Password: "password",
			},
			prepare: func(u model.User) {
				mock.ExpectQuery(`INSERT INTO converter."user" (.+) RETURNING`).
					WithArgs(u.Username, u.Password).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(u.ID))
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(tc.expectedUser)

			userID, err := repo.InsertUser(tc.expectedUser.Username, tc.expectedUser.Password)
			assertNoError(t, fmt.Errorf("error when inserting user: '%w'", err))

			if tc.expectedUser.ID != userID {
				t.Errorf("expected user id to be %q, but got %q", tc.expectedUser.ID, userID)
			}

			err = mock.ExpectationsWereMet()
			assertNoError(t, fmt.Errorf("there were unfulfilled expectations: %w", err))
		})
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if errors.Unwrap(err) != nil {
		t.Errorf("No error expected, got %v ", err)
	}
}
