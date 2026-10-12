package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	TooShort error = errors.New("Password is too short")
)

func HashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", TooShort
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 3)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}
