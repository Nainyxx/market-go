package models

import (
	"errors"
	"fmt"

	"github.com/Nainyxx/market-go/services/authservice/utils"
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`

	Phone        string `json:"phone"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

var (
	ErrInvalidData  = errors.New("input data is invalid")
	ErrPasswordHash = errors.New("failed to hash password")
	ErrAlreadyExist = errors.New("user already exist")
	ErrNotFound     = errors.New("user not found")
)

type CreateError struct {
	Field string
	Err   error
}

func (e *CreateError) Error() string {
	return fmt.Sprintf("failed to create user due to invalid %s: %v", e.Field, e.Err)
}

func (e *CreateError) Unwrap() error {
	return e.Err
}

func CreateUser(phone, email, password string) (*User, error) {
	if err := utils.IsPhoneValid(phone); err != nil {
		return nil, &CreateError{Field: "phone", Err: fmt.Errorf("%w: %w", ErrInvalidData, err)}
	}
	if err := utils.IsEmailValid(email); err != nil {
		return nil, &CreateError{Field: "email", Err: fmt.Errorf("%w: %w", ErrInvalidData, err)}
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, &CreateError{Field: "password", Err: fmt.Errorf("%w: %w", ErrInvalidData, err)}
	}

	return &User{
		ID:           uuid.New(),
		Phone:        phone,
		Email:        email,
		PasswordHash: hash,
	}, nil
}
