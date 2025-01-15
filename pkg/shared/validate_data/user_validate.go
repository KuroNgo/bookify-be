package validate_data

import (
	"bookify/internal/domain"
	"bookify/pkg/shared/helper"
	"errors"
)

func ValidateUser(input *domain.InputUser) error {
	if input.PasswordHash == "" {
		return errors.New("the user's information cannot be empty")
	}

	if input.Email == "" {
		return errors.New("the user's information cannot be empty")
	}

	if !helper.EmailValid(input.Email) {
		return errors.New("email Invalid ")
	}

	if input.FullName == "" {
		return errors.New("the user's information cannot be empty")
	}

	return nil
}

func ValidateUser2(input *domain.SignupUser) error {
	if input.Password == "" {
		return errors.New("the user's information cannot be empty")
	}

	if input.Email == "" {
		return errors.New("the user's information cannot be empty")
	}

	if !helper.EmailValid(input.Email) {
		return errors.New("email Invalid ")
	}

	return nil
}

func ValidateUser3(input *domain.User) error {
	if input.PasswordHash == "" {
		return errors.New("the user's information cannot be empty")
	}

	if input.Email == "" {
		return errors.New("the user's information cannot be empty")
	}

	if !helper.EmailValid(input.Email) {
		return errors.New("email Invalid ")
	}

	return nil
}
