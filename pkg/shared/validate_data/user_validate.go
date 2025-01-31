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

func ValidateUser4(input *domain.UpdateUserSettings) error {
	if input.Gender == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Vocation == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Address == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.City == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Region == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.DateOfBirth == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.FullName == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	return nil
}
