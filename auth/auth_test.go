package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		shouldError bool
	}{
		{
			name:        "Valid input",
			password:    "password",
			shouldError: false,
		},
		{
			name:        "Empty input",
			password:    "",
			shouldError: false,
		},
		{
			name:        "Input just short enough",
			password:    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			shouldError: false,
		},
		{
			name:        "Input too long",
			password:    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			shouldError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := hashPassword(test.password)
			if test.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	password := "password"
	hash, _ := hashPassword(password)

	assert.True(t, ValidatePassword(password, hash))
	assert.False(t, ValidatePassword("invalid", hash))
}
