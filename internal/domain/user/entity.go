// Package user implements a database user, his validation, queries and handlers.
package user

import "errors"

// User represents a database user.
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Created  string `json:"created"`
	Updated  string `json:"updated"`
}

func (u *User) validateCredentials() error {
	if u.Username == "" {
		return errors.New("username is missing")
	}
	if u.Password == "" {
		return errors.New("password is missing")
	}
	// return  errors.New("invalid login/password")
	return nil
}

func (u *User) exists() error {
	// return  errors.New("the user already exists")
	return nil
}

// Validate validates user's data.
func (u *User) validate() error {
	if err := u.validateCredentials(); err != nil {
		return err
	}

	if err := u.exists(); err != nil {
		return err
	}

	return nil
}
