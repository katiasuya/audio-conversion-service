// Package hash provide functions to hash and compare passwords.
package hash

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordHash hashes user's password.
func PasswordHash(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	return string(hash), err
}

// CheckPasswordHash checks whether the hashed passwords match.
func CheckPasswordHash(pwd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil
}
