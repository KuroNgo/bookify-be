package validate_data

import (
	"bookify/internal/domain"
	"bookify/pkg/shared/constants"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateEventWishlist(input domain.EventWishlist) error {
	if input.ID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}
	if input.EventID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}
	if input.UserID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}

	return nil
}

func ValidateEventWishlistInput(input *domain.EventWishlistInput) error {
	if input.EventID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}
	if input.UserID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}

	return nil
}
