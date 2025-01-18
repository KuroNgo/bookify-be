package validate_data

import (
	"bookify/internal/domain"
	"bookify/pkg/shared/helper"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateOrganization(input *domain.Organization) error {
	if input.ID == primitive.NilObjectID {
		return errors.New("the partner's information cannot be empty")
	}

	if input.Name == "" {
		return errors.New("the partner's information cannot be empty")
	}

	if input.Email == "" {
		return errors.New("the partner's information cannot be empty")
	}
	if !helper.EmailValid(input.Email) {
		return errors.New("email Invalid ")
	}

	if input.Phone == "" {
		return errors.New("the partner's information cannot be empty")
	}
	if !helper.PhoneValid(input.Phone) {
		return errors.New("phone Invalid")
	}

	return nil
}

func ValidateOrganizationInput(input *domain.OrganizationInput) error {
	if input.Name == "" {
		return errors.New("the partner's information cannot be empty")
	}

	if input.Email == "" {
		return errors.New("the partner's information cannot be empty")
	}
	if !helper.EmailValid(input.Email) {
		return errors.New("email Invalid ")
	}

	if input.Phone == "" {
		return errors.New("the partner's information cannot be empty")
	}
	if !helper.PhoneValid(input.Phone) {
		return errors.New("phone Invalid")
	}

	return nil
}
