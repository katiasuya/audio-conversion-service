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

var formats = map[string]string{"MP3": "audio/mpeg", "WAV": "audio/wave"}

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
		return errors.New("invalid target format: need mp3 or wav")
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
