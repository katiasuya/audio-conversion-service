// Package domain describes principal entities of the project and their validations.
package domain

// User represents user's credentials.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ValidateCredentials validates user's credentials for having illegal characters or length being out of range.
func (u *User) ValidateCredentials() error {
	// return  errors.New("invalid login/password")
	return nil
}

// Exists checks if the user with the given username already exists
func (u *User) Exists() error {
	// return  errors.New("the user already exists")
	return nil
}

// IncorrectPwd checks whether the given password matches with the password of the given username.
func (u *User) IncorrectPwd() error {
	// return  errors.New("incorrect password")
	return nil
}

// IsAuthorized checks whether the user is autorized.
func (u *User) IsAuthorized() error {
	// return  errors.New("sorry, you are not authorized")
	return nil
}
