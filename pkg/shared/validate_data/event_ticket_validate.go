package validate_data

import (
	"bookify/internal/domain"
	"bookify/pkg/shared/constants"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateEventTicket(input domain.EventTicket) error {
	if input.ID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}
	if input.EventID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Price < 0 {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Quantity < 0 {
		return errors.New(constants.MsgInvalidInput)
	}

	return nil
}

func ValidateEventTicketInput(input *domain.EventTicketInput) error {
	if input.EventID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Price < 0 {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Quantity < 0 {
		return errors.New(constants.MsgInvalidInput)
	}

	return nil
}
