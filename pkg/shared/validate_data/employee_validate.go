package validate_data

import (
	"bookify/internal/domain"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/helper"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateEmployee(input *domain.Employee) error {
	if input.ID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.OrganizationID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Email == "" {
		return errors.New(constants.MsgInvalidInput)
	}
	if !helper.EmailValid(input.Email) {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.FirstName == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.LastName == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.JobTitle == "" {
		return errors.New(constants.MsgInvalidInput)
	}
	return nil
}

func ValidateEmployeeInput(input *domain.EmployeeInput) error {
	if input.OrganizationID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Email == "" {
		return errors.New(constants.MsgInvalidInput)
	}
	if !helper.EmailValid(input.Email) {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.FirstName == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.LastName == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.JobTitle == "" {
		return errors.New(constants.MsgInvalidInput)
	}
	return nil
}
