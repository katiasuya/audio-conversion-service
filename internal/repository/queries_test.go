package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

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
				Name:     "Yesterday",
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

func TestGetIDAndPasswordByUsername(t *testing.T) {
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
				Username: "user123",
				Password: "qwerty123",
			},
			prepare: func(u model.User) {
				mock.ExpectQuery(`SELECT (.*) FROM converter."user" WHERE username=?`).
					WithArgs(u.Username).
					WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).
						AddRow(u.ID, u.Password))
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(tc.expectedUser)

			gotID, gotPassword, err := repo.GetIDAndPasswordByUsername(tc.expectedUser.Username)
			assertNoError(t, fmt.Errorf("error when getting id and password: '%w'", err))

			if tc.expectedUser.ID != gotID || tc.expectedUser.Password != gotPassword {
				t.Errorf("expected user to have id=%s and password=%s, but got id=%s and password=%s",
					tc.expectedUser.ID, tc.expectedUser.Password, gotID, gotPassword)
			}

			err = mock.ExpectationsWereMet()
			assertNoError(t, fmt.Errorf("there were unfulfilled expectations: %w", err))
		})
	}
}

func TestGetRequestHistory(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}
	defer repo.Close()

	cases := []struct {
		name            string
		expectedRequest model.Request
		prepare         func(model.Request)
	}{
		{
			name: "success",
			expectedRequest: model.Request{
				ID:           "1",
				AudioName:    "Yesterday",
				SourceFormat: "mp3",
				TargetFormat: "wav",
				Created:      time.Time{},
				Updated:      time.Time{},
				Status:       "queued",
			},
			prepare: func(r model.Request) {
				mock.ExpectQuery(`SELECT r.id, a.name, r.source_format, r.target_format, r.created, r.updated, r.status
				FROM converter.request r JOIN converter.audio a ON a.id = r.source_id WHERE r.user_id=?`).
					WithArgs(r.UserID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "source_format", "target_format", "created", "updated", "status"}).
						AddRow(r.ID, r.AudioName, r.SourceFormat, r.TargetFormat, r.Created, r.Updated, r.Status))
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(tc.expectedRequest)

			gotRequest, err := repo.GetRequestHistory(tc.expectedRequest.UserID)
			assertNoError(t, fmt.Errorf("error when getting audio by ID: '%w'", err))

			if tc.expectedRequest != gotRequest[0] {
				t.Errorf("expected user to be %+v, but got %+v", tc.expectedRequest, gotRequest)
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
				Username: "User123",
				Password: "qwerty123",
			},
			prepare: func(u model.User) {
				mock.ExpectQuery(`INSERT INTO converter."user" (.*) RETURNING`).
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

func TestInsertAudio(t *testing.T) {
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
				Name:     "Yesterday",
				Format:   "mp3",
				Location: "location",
			},
			prepare: func(a model.Audio) {
				mock.ExpectQuery(`INSERT INTO converter.audio (.+) RETURNING`).
					WithArgs(a.Name, a.Format, a.Location).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(a.ID))
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(tc.expectedAudio)

			audioID, err := repo.InsertAudio(tc.expectedAudio.Name, tc.expectedAudio.Format, tc.expectedAudio.Location)
			assertNoError(t, fmt.Errorf("error when inserting audio: '%w'", err))

			if tc.expectedAudio.ID != audioID {
				t.Errorf("expected audio id to be %q, but got %q", tc.expectedAudio.ID, audioID)
			}

			err = mock.ExpectationsWereMet()
			assertNoError(t, fmt.Errorf("there were unfulfilled expectations: %w", err))
		})
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if errors.Unwrap(err) != nil {
		t.Errorf("No error expected, got '%v' ", err)
	}
}
