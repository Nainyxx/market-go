package utils

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrValidation = errors.New("validation failed")

func validationErr(msg string) error {
	return fmt.Errorf("%w: %s", ErrValidation, msg)
}

func IsEmailValid(email string) error {
	if email == "" {
		return validationErr("email cannot be empty")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return validationErr("invalid email format")
	}

	return nil
}

func IsNameValid(name string) error {
	if name == "" || utf8.RuneCountInString(name) > 64 {
		return validationErr("invalid name")
	}
	if ContainsOnlyLetter(name) != nil {
		return validationErr("invalid name")
	}
	return nil
}

func IsSurnameValid(surname string) error {
	if utf8.RuneCountInString(surname) > 64 {
		return validationErr("invalid surname")
	}
	if ContainsOnlyLetter(surname) != nil {
		return validationErr("invalid surname")
	}
	return nil
}

func ContainsOnlyLetter(str string) error {
	for _, v := range str {
		if !unicode.IsLetter(v) {
			return validationErr("invalid name")
		}
	}
	return nil
}

func ContainsOnlyNumber(str string) error {
	for _, v := range str {
		if !unicode.IsNumber(v) {
			return validationErr("invalid input")
		}
	}
	return nil
}

func IsLoginValid(login string) error {
	if login == "" || utf8.RuneCountInString(login) > 32 {
		return validationErr("invalid login")
	}
	for _, v := range login {
		if !(unicode.IsLetter(v) || unicode.IsNumber(v) || string(v) == "_" || string(v) == "-") {
			return validationErr("invalid login")
		}
	}
	return nil
}

func IsPhoneValid(phone string) error {
	if phone == "" {
		return validationErr("phone cannot be empty")
	}

	digitsOnly := phone
	if phone[0] == '+' {
		digitsOnly = phone[1:]
	}

	if ContainsOnlyNumber(digitsOnly) != nil {
		return validationErr("invalid phone")
	}

	if len(digitsOnly) < 10 || len(digitsOnly) > 15 {
		return validationErr("phone length is invalid")
	}

	return nil
}

func IsImageURLValid(imageURL string) error {
	if imageURL == "" {
		return nil
	}

	u, err := url.ParseRequestURI(imageURL)
	if err != nil {
		return validationErr("invalid URL format")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return validationErr("URL must use http or https scheme")
	}

	path := strings.ToLower(u.Path)
	isImage := false

	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, ext := range validExtensions {
		if strings.HasSuffix(path, ext) {
			isImage = true
			break
		}
	}

	if !isImage {
		return validationErr("URL must point to a valid image file (.jpg, .png, .webp, etc.)")
	}

	return nil
}

func IsDescriptionValid(desc string) error {
	if utf8.RuneCountInString(desc) > 255 {
		return validationErr("description is too long")
	}
	return nil
}
