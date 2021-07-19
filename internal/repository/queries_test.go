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
	"github.com/lib/pq"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(context.Background(),
			fmt.Errorf("an error '%s' was not expected when opening a stub database connection", err))
	}

	return db, mock
}

func NewMockWithMatcher() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		logger.Fatal(context.Background(),
			fmt.Errorf("an error '%s' was not expected when opening a stub database connection", err))
	}

	return db, mock
}

func TestGetIDAndPasswordByUsername(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}
	defer repo.Close()

	var dbErr = errors.New("test db error")

	cases := []struct {
		name          string
		expectedUser  model.User
		expectedError error
		prepare       func(model.User)
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
		{
			name:          "no such user",
			expectedError: ErrNoSuchUser,
			prepare: func(u model.User) {
				mock.ExpectQuery(`SELECT (.*) FROM converter."user" WHERE username=?`).
					WithArgs(u.Username).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:          "db call error",
			expectedError: dbErr,
			prepare: func(u model.User) {
				mock.ExpectQuery(`SELECT (.*) FROM converter."user" WHERE username=?`).
					WithArgs(u.Username).
					WillReturnError(dbErr)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(tc.expectedUser)

			gotID, gotPassword, err := repo.GetIDAndPasswordByUsername(tc.expectedUser.Username)
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error to be %v, but got %v", tc.expectedError, err)
			}

			if tc.expectedUser.ID != gotID || tc.expectedUser.Password != gotPassword {
				t.Errorf("expected user to have id=%s and password=%s, but got id=%s and password=%s",
					tc.expectedUser.ID, tc.expectedUser.Password, gotID, gotPassword)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestGetRequestHistory(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}
	defer repo.Close()

	var dbErr = errors.New("test db error")

	cases := []struct {
		name            string
		expectedHistory model.Request
		expectedError   error
		prepare         func(model.Request)
	}{
		{
			name: "success",
			expectedHistory: model.Request{
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
		{
			name:          "db call error",
			expectedError: dbErr,
			prepare: func(r model.Request) {
				mock.ExpectQuery(`SELECT r.id, a.name, r.source_format, r.target_format, r.created, r.updated, r.status
				FROM converter.request r JOIN converter.audio a ON a.id = r.source_id WHERE r.user_id=?`).
					WithArgs(r.UserID).
					WillReturnError(dbErr)
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(tc.expectedHistory)

			gotRequest, err := repo.GetRequestHistory(tc.expectedHistory.UserID)
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error to be %v, but got %v", tc.expectedError, err)
			}

			if len(gotRequest) > 0 {
				if tc.expectedHistory != gotRequest[0] {
					t.Errorf("expected user to be %+v, but got %+v", tc.expectedHistory, gotRequest)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestGetAudioByID(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}
	defer repo.Close()

	var dbErr = errors.New("test db error")

	cases := []struct {
		name          string
		expectedAudio model.Audio
		expectedError error
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
				mock.ExpectQuery("SELECT (.*) FROM converter.audio WHERE id=?").
					WithArgs(a.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "format", "location"}).
						AddRow(a.ID, a.Name, a.Format, a.Location))
			},
		},
		{
			name:          "no such audio",
			expectedError: ErrNoSuchAudio,
			prepare: func(a model.Audio) {
				mock.ExpectQuery("SELECT (.*) FROM converter.audio WHERE id=?").
					WithArgs(a.ID).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:          "db call error",
			expectedError: dbErr,
			prepare: func(a model.Audio) {
				mock.ExpectQuery("SELECT (.*) FROM converter.audio WHERE id=?").
					WithArgs(a.ID).
					WillReturnError(dbErr)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(tc.expectedAudio)

			gotAudio, err := repo.GetAudioByID(tc.expectedAudio.ID)
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error to be %v, but got %v", tc.expectedError, err)
			}

			if gotAudio != tc.expectedAudio {
				t.Errorf("expected audio to be %+v, but got %+v", tc.expectedAudio, gotAudio)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestInsertUser(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}
	defer repo.Close()

	var dbErr = errors.New("test db error")

	cases := []struct {
		name          string
		expectedUser  model.User
		expectedError error
		prepare       func(model.User)
	}{
		{
			name: "success",
			expectedUser: model.User{
				ID:       "1",
				Username: "user123",
				Password: "qwerty123",
			},
			prepare: func(u model.User) {
				mock.ExpectQuery(`INSERT INTO converter."user" (.*) RETURNING id`).
					WithArgs(u.Username, u.Password).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(u.ID))
			},
		},
		{
			name:          "user already exists",
			expectedError: ErrUserAlreadyExists,
			prepare: func(u model.User) {
				mock.ExpectQuery(`INSERT INTO converter."user" (.*) RETURNING id`).
					WithArgs(u.Username, u.Password).
					WillReturnError(&pq.Error{Code: codeUniqueViolation})
			},
		},
		{
			name:          "db call error",
			expectedError: dbErr,
			prepare: func(u model.User) {
				mock.ExpectQuery(`INSERT INTO converter."user" (.*) RETURNING id`).
					WithArgs(u.Username, u.Password).
					WillReturnError(dbErr)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(tc.expectedUser)

			userID, err := repo.InsertUser(tc.expectedUser.Username, tc.expectedUser.Password)
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error to be %v, but got %v", tc.expectedError, err)
			}

			if tc.expectedUser.ID != userID {
				t.Errorf("expected user id to be %q, but got %q", tc.expectedUser.ID, userID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestInsertAudio(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}
	defer repo.Close()

	var dbErr = errors.New("test db error")

	cases := []struct {
		name          string
		expectedAudio model.Audio
		expectedError error
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
				mock.ExpectQuery(`INSERT INTO converter.audio (.+) RETURNING id`).
					WithArgs(a.Name, a.Format, a.Location).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(a.ID))
			},
		},
		{
			name:          "db call error",
			expectedError: dbErr,
			prepare: func(a model.Audio) {
				mock.ExpectQuery(`INSERT INTO converter.audio (.+) RETURNING id`).
					WithArgs(a.Name, a.Format, a.Location).
					WillReturnError(dbErr)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(tc.expectedAudio)

			audioID, err := repo.InsertAudio(tc.expectedAudio.Name, tc.expectedAudio.Format, tc.expectedAudio.Location)
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error to be %v, but got %v", tc.expectedError, err)
			}

			if tc.expectedAudio.ID != audioID {
				t.Errorf("expected audio id to be %q, but got %q", tc.expectedAudio.ID, audioID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateRequest(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}
	defer repo.Close()

	var dbErr = errors.New("test db error")
	var rowsErr = errors.New("test rowsAffected error")

	cases := []struct {
		name          string
		args          model.Request
		expectedError error
		prepare       func(model.Request)
	}{
		{
			name: "success",
			args: model.Request{
				ID:       "1",
				TargetID: "2",
				Status:   "queued",
			},
			prepare: func(r model.Request) {
				mock.ExpectExec("UPDATE converter.request SET (.*) WHERE id=?").
					WithArgs(r.ID, r.TargetID, r.Status).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "empty target id",
			args: model.Request{
				ID:       "1",
				TargetID: "",
				Status:   "queued",
			},
			prepare: func(r model.Request) {
				mock.ExpectExec("UPDATE converter.request SET (.*) WHERE id=?").
					WithArgs(r.ID, sql.NullString{}, r.Status).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:          "no such request",
			expectedError: ErrNoSuchRequest,
			prepare: func(r model.Request) {
				mock.ExpectExec("UPDATE converter.request SET (.*) WHERE id=?").
					WithArgs(r.ID, sql.NullString{}, r.Status).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
		{
			name:          "RowsAffected() error",
			expectedError: rowsErr,
			prepare: func(r model.Request) {
				mock.ExpectExec("UPDATE converter.request SET (.*) WHERE id=?").
					WithArgs(r.ID, sql.NullString{}, r.Status).
					WillReturnResult(sqlmock.NewErrorResult(rowsErr))
			},
		},
		{
			name:          "db call error",
			expectedError: dbErr,
			prepare: func(r model.Request) {
				mock.ExpectExec("UPDATE converter.request SET (.*) WHERE id=?").
					WithArgs(r.ID, sql.NullString{}, r.Status).
					WillReturnError(dbErr)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(tc.args)

			err := repo.UpdateRequest(tc.args.ID, tc.args.Status, tc.args.TargetID)
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error to be %v, but got %v", tc.expectedError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
