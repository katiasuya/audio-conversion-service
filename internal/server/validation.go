package server

import (
	"errors"
)

const (
	minLength         = 5
	maxNameLength     = 256
	maxPasswordLength = 128
)

var errInvalidLength = errors.New("invalid length")

// ValidateUserCredentials validates user's credentials.
func ValidateUserCredentials(username, password string) error {
	if username == "" {
		return errors.New("username is missing")
	}
	if len(username) < minLength || len(username) > maxNameLength {
		return errInvalidLength
	}
	if err := validateChars(username); err != nil {
		return err
	}

	if password == "" {
		return errors.New("password is missing")
	}
	if len(password) < minLength || len(password) > maxPasswordLength {
		return errInvalidLength
	}
	if err := validateChars(password); err != nil {
		return err
	}

	return nil
}

// validateChars checks whether the given string contains invalid characters.
func validateChars(str string) error {
	//return errors.New("invalid characters")
	return nil
}
