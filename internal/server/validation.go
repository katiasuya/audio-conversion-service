package server

import (
	"errors"
	"fmt"
	"strings"
)

const (
	minLength = 6
	maxLength = 128
)

var formats = map[string]string{"mp3": "audio/mpeg", "wav": "audio/wave"}

// ValidateUserCredentials validates user's credentials.
func ValidateUserCredentials(username, password string) error {
	if username == "" {
		return errors.New("username is missing")
	}
	if len(username) < minLength || len(username) > maxLength {
		return fmt.Errorf("invalid length, username must be from %d to %d characters", minLength, maxLength)
	}

	if err := validateChars(username); err != nil {
		return err
	}

	if password == "" {
		return errors.New("password is missing")
	}
	if len(password) < minLength || len(password) > maxLength {
		return fmt.Errorf("invalid length, password must be from %d to %d characters", minLength, maxLength)
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
		return errors.New("invalid target format, need mp3 or wav")
	}

	if err := validateChars(name); err != nil {
		return err
	}

	return nil
}

// validateChars checks whether the given string contains invalid characters.
func validateChars(str string) error {
	const invalidChars = `:;<>\{}[]+=?&," `
	if strings.ContainsAny(str, invalidChars) {
		return fmt.Errorf("invalid character(s), you can't use %sand space", invalidChars)
	}

	return nil
}
