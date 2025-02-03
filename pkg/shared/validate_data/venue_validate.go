package validate_data

import (
	"bookify/internal/domain"
	"bookify/pkg/shared/constants"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateVenue(input *domain.Venue) error {
	if input.ID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.AddressLine == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.City == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	//if input.State == "" {
	//	return errors.New(constants.MsgInvalidInput)
	//}

	if input.Country == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	//if input.PostalCode == "" {
	//	return errors.New(constants.MsgInvalidInput)
	//}

	return nil
}

func ValidateVenueInput(input *domain.VenueInput) error {
	if input.AddressLine == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.City == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	//if input.State == "" {
	//	return errors.New(constants.MsgInvalidInput)
	//}

	if input.Country == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	//if input.PostalCode == "" {
	//	return errors.New(constants.MsgInvalidInput)
	//}

	return nil
}
