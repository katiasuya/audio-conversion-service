package validation

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
			exp:      ErrMissingUsername,
		},
		{
			name:     "missing password",
			username: "jonathan",
			password: "",
			exp:      ErrMissingPassword,
		},
		{
			name:     "invalid username length",
			username: "john",
			password: "1",
			exp:      ErrInvalidLength,
		},
		{
			name:     "invalid password length",
			username: "jonathan",
			password: "123",
			exp:      ErrInvalidLength,
		},
		{
			name:     "invalid username chars",
			username: "joh n123",
			password: "1",
			exp:      ErrInvalidChars,
		},
		{
			name:     "invalid password chars",
			username: "jonathan",
			password: "qwerty&123",
			exp:      ErrInvalidChars,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ValidateUserCredentials(tt.username, tt.password)
			if res != tt.exp {
				t.Errorf("Expected %t, got %t", tt.exp, res)
			}
		})
	}
}
