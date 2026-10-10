package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", errors.New("Password is too short")
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 3)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}
