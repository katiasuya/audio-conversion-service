package server

import "testing"

// TestValidateUserCredentials tests ValidateUserCredentials function.
func TestValidateUserCredentials(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		exp      error
	}{
		{
			name:     "missing username",
			username: "",
			password: "qwerty123",
			exp:      errMissingUsername,
		},
		{
			name:     "missing password",
			username: "jonathan",
			password: "",
			exp:      errMissingPassword,
		},
		{
			name:     "invalid username length",
			username: "john",
			password: "1",
			exp:      errInvalidLength,
		},
		{
			name:     "invalid password length",
			username: "jonathan",
			password: "123",
			exp:      errInvalidLength,
		},
		{
			name:     "invalid username chars",
			username: "joh n123",
			password: "1",
			exp:      errInvalidChars,
		},
		{
			name:     "invalid password chars",
			username: "jonathan",
			password: "qwerty&123",
			exp:      errInvalidChars,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := UserCredentials(tt.username, tt.password)
			if res != tt.exp {
				t.Errorf("Expected %t, got %t", tt.exp, res)
			}
		})
	}
}
