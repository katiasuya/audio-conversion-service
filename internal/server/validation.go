package server

import (
	"errors"
	"fmt"
)

// ValidateUserCredentials validates user's credentials.
func ValidateUserCredentials(username, password string) error {
	if username == "" {
		return errors.New("username is missing")
	}
	err := validate(username)
	if err != nil {
		return err
	}

	if password == "" {
		return errors.New("password is missing")
	}
	err = validate(password)
	if err != nil {
		return err
	}

	return nil
}

func validate(str string) error {
	const minLength, maxLength = 5, 15
	if len(str) < minLength || len(str) > maxLength {
		lenErr := fmt.Sprintf("%s: invalid length. It must be from 5 to 15", str)
		return errors.New(lenErr)
	}

	//return errors.New("invalid characters")

	return nil
}
