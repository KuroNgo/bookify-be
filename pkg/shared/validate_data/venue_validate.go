package validate_data

import (
	"bookify/internal/domain"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateVenue(input *domain.Venue) error {
	if input.ID == primitive.NilObjectID {
		return errors.New("the venue's information cannot be empty")
	}

	if input.AddressLine == "" {
		return errors.New("the venue's information cannot be empty")
	}

	if input.City == "" {
		return errors.New("the venue's information cannot be empty")
	}
	if input.State == "" {
		return errors.New("the venue's information cannot be empty")
	}

	if input.Country == "" {
		return errors.New("the venue's information cannot be empty")
	}

	if input.PostalCode == "" {
		return errors.New("the venue's information cannot be empty")
	}

	return nil
}

func ValidateVenueInput(input *domain.VenueInput) error {
	if input.AddressLine == "" {
		return errors.New("the venue's information cannot be empty")
	}

	if input.City == "" {
		return errors.New("the venue's information cannot be empty")
	}
	if input.State == "" {
		return errors.New("the venue's information cannot be empty")
	}

	if input.Country == "" {
		return errors.New("the venue's information cannot be empty")
	}

	if input.PostalCode == "" {
		return errors.New("the venue's information cannot be empty")
	}

	return nil
}
