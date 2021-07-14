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

var (
	audio = model.Audio{
		ID:       "1",
		Name:     "Yesterday",
		Format:   "mp3",
		Location: "location",
	}
	user = model.User{
		ID:       "1",
		Username: "user123",
		Password: "qwerty123",
	}
	request = model.Request{
		ID:           "1",
		AudioName:    "Yesterday",
		SourceFormat: "mp3",
		TargetFormat: "wav",
		Created:      time.Time{},
		Updated:      time.Time{},
		Status:       "queued",
	}
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
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

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT (.*) FROM converter."user" WHERE username=?`).
			WithArgs(user.Username).
			WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).
				AddRow(user.ID, user.Password))

		gotID, gotPassword, err := repo.GetIDAndPasswordByUsername(user.Username)
		assertNoError(t, fmt.Errorf("error when getting id and password: '%w'", err))

		if user.ID != gotID || user.Password != gotPassword {
			t.Errorf("expected user to have id=%s and password=%s, but got id=%s and password=%s",
				user.ID, user.Password, gotID, gotPassword)
		}

		err = mock.ExpectationsWereMet()
		assertNoError(t, fmt.Errorf("there were unfulfilled expectations: %w", err))
	})

	t.Run("no such user", func(t *testing.T) {
		mock.ExpectQuery(`SELECT (.*) FROM converter."user" WHERE username=?`).
			WithArgs(user.Username).
			WillReturnRows(sqlmock.NewRows([]string{"id", "password"}))

		_, _, err := repo.GetIDAndPasswordByUsername(user.Username)
		assertError(t, fmt.Errorf("error when getting id and password: '%w'", err), ErrNoSuchUser)

		err = mock.ExpectationsWereMet()
		assertNoError(t, fmt.Errorf("there were unfulfilled expectations: %w", err))
	})
}

func TestGetRequestHistory(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}
	defer repo.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT r.id, a.name, r.source_format, r.target_format, r.created, r.updated, r.status
		FROM converter.request r JOIN converter.audio a ON a.id = r.source_id WHERE r.user_id=?`).
			WithArgs(request.UserID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "source_format", "target_format", "created", "updated", "status"}).
				AddRow(request.ID, request.AudioName, request.SourceFormat, request.TargetFormat, request.Created, request.Updated, request.Status))

		gotRequest, err := repo.GetRequestHistory(request.UserID)
		assertNoError(t, fmt.Errorf("error when getting audio by ID: '%w'", err))

		if request != gotRequest[0] {
			t.Errorf("expected user to be %+v, but got %+v", request, gotRequest)
		}

		err = mock.ExpectationsWereMet()
		assertNoError(t, fmt.Errorf("there were unfulfilled expectations: %w", err))
	})

}
func TestGetAudioByID(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}
	defer repo.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.*) FROM converter.audio WHERE id=?").
			WithArgs(audio.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "format", "location"}).
				AddRow(audio.ID, audio.Name, audio.Format, audio.Location))

		gotAudio, err := repo.GetAudioByID(audio.ID)
		assertNoError(t, fmt.Errorf("error when getting audio by ID: '%w'", err))

		if gotAudio != audio {
			t.Errorf("expected user to be %+v, but got %+v", audio, gotAudio)
		}

		err = mock.ExpectationsWereMet()
		assertNoError(t, fmt.Errorf("there were unfulfilled expectations: %w", err))
	})

	t.Run("no such audio", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.*) FROM converter.audio WHERE id=?").
			WithArgs(audio.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "format", "location"}))

		_, err := repo.GetAudioByID(audio.ID)
		assertError(t, fmt.Errorf("error when getting audio by ID: '%w'", err), ErrNoSuchAudio)

		err = mock.ExpectationsWereMet()
		assertNoError(t, fmt.Errorf("there were unfulfilled expectations: %w", err))
	})
}

func TestInsertUser(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}
	defer repo.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO converter."user" (.*) RETURNING`).
			WithArgs(user.Username, user.Password).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(user.ID))

		userID, err := repo.InsertUser(user.Username, user.Password)
		assertNoError(t, fmt.Errorf("error when inserting user: '%w'", err))

		if user.ID != userID {
			t.Errorf("expected user id to be %q, but got %q", user.ID, userID)
		}

		err = mock.ExpectationsWereMet()
		assertNoError(t, fmt.Errorf("there were unfulfilled expectations: %w", err))
	})
}

func TestInsertAudio(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}
	defer repo.Close()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO converter.audio (.+) RETURNING`).
			WithArgs(audio.Name, audio.Format, audio.Location).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(audio.ID))

		audioID, err := repo.InsertAudio(audio.Name, audio.Format, audio.Location)
		assertNoError(t, fmt.Errorf("error when inserting audio: '%w'", err))

		if audio.ID != audioID {
			t.Errorf("expected audio id to be %q, but got %q", audio.ID, audioID)
		}

		err = mock.ExpectationsWereMet()
		assertNoError(t, fmt.Errorf("there were unfulfilled expectations: %w", err))
	})
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if errors.Unwrap(err) != nil {
		t.Errorf("No error expected, got '%v' ", err)
	}
}

func assertError(t *testing.T, err, assertErr error) {
	t.Helper()
	if errors.Unwrap(err) != assertErr {
		t.Errorf("Expected '%v', got '%v' error", err, assertErr)
	}
}
