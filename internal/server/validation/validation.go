package validation

import (
	"errors"
	"fmt"
	"strings"
)

const (
	minLength = 6
	maxLength = 128
)
const invalidChars = `:;<>\{}[]+=?&," `

var formats = map[string]string{"mp3": "audio/mpeg", "wav": "audio/wave"}

var (
	ErrMissingUsername = errors.New("username is missing")
	ErrMissingPassword = errors.New("password is missing")
	ErrInvalidLength   = fmt.Errorf("invalid length: username and password must be from %d to %d characters", minLength, maxLength)
	ErrInvalidChars    = fmt.Errorf("invalid character(s): you can't use %sor space character(s)", invalidChars)
)

// ValidateUserCredentials validates user's credentials.
func ValidateUserCredentials(username, password string) error {
	if username == "" {
		return ErrMissingUsername
	}
	if len(username) < minLength || len(username) > maxLength {
		return ErrInvalidLength
	}

	if containsInvalidChars(username) {
		return ErrInvalidChars
	}

	if password == "" {
		return ErrMissingPassword
	}
	if len(password) < minLength || len(password) > maxLength {
		return ErrInvalidLength
	}
	if containsInvalidChars(password) {
		return ErrInvalidChars
	}

	return nil
}

// ValidateRequest validates conversion request body.
func ValidateRequest(name, sourceFormat, targetFormat, sourceContentType string) error {
	if sourceFormat == "" {
		return errors.New("source format is missing")
	}
	if targetFormat == "" {
		return errors.New("target format is missing")
	}
	if formats[sourceFormat] != sourceContentType {
		return errors.New("wrong source format for the file")
	}
	if sourceFormat == targetFormat {
		return errors.New("source and target formats can't be equal")
	}
	if _, ok := formats[targetFormat]; !ok {
		return errors.New("invalid target format, need mp3 or wav")
	}

	if containsInvalidChars(name) {
		return ErrInvalidChars
	}

	return nil
}

// containsInvalidChars checks whether the given string contains invalid characters.
func containsInvalidChars(str string) bool {
	return strings.ContainsAny(str, invalidChars)
}
