package validate_data

import (
	"bookify/internal/domain"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/helper"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidatePartner(input *domain.Partner) error {
	if input.ID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Name == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Email == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if !helper.EmailValid(input.Email) {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Phone == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if !helper.PhoneValid(input.Phone) {
		return errors.New(constants.MsgInvalidInput)
	}

	return nil
}

func ValidatePartnerInput(input *domain.PartnerInput) error {
	if input.Name == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Email == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if !helper.EmailValid(input.Email) {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Phone == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if !helper.PhoneValid(input.Phone) {
		return errors.New(constants.MsgInvalidInput)
	}

	return nil
}
