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

// Formats lists all possible audio formats.
var formats = []string{"MP3", "WAV"}

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

// ValidateRequest validates conversion request body.
func ValidateRequest(name, sourceFormat, targetFormat string) error {
	if sourceFormat == "" {
		return errors.New("source format is missing")
	}
	if targetFormat == "" {
		return errors.New("target format is missing")
	}

	if !contains(sourceFormat, formats) {
		return errors.New("invalid source format: need MP3 or WAV")
	}
	if !contains(targetFormat, formats) {
		return errors.New("invalid target format: need MP3 or WAV")
	}

	if sourceFormat == targetFormat {
		return errors.New("source and target formats can't be equal")
	}

	if err := validateChars(name); err != nil {
		return err
	}

	return nil
}

// validateChars checks whether the given string contains invalid characters.
func validateChars(str string) error {
	//return errors.New("invalid characters")
	return nil
}

// validateChars checks whether the given format is present in formats slice.
func contains(format string, formats []string) bool {
	for _, f := range formats {
		if format == f {
			return true
		}
	}
	return false
}
