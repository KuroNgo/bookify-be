package validate_data

import (
	"bookify/internal/domain"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/helper"
	"errors"
)

func ValidateUser(input *domain.InputUser) error {
	if input.PasswordHash == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Email == "" {
		return errors.New(constants.MsgInvalidInput)

	}

	if !helper.EmailValid(input.Email) {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.FullName == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	return nil
}

func ValidateUser2(input *domain.SignupUser) error {
	if input.Password == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Email == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if !helper.EmailValid(input.Email) {
		return errors.New(constants.MsgInvalidInput)
	}

	return nil
}

func ValidateUser3(input *domain.User) error {
	if input.PasswordHash == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Email == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if !helper.EmailValid(input.Email) {
		return errors.New(constants.MsgInvalidInput)
	}

	return nil
}
