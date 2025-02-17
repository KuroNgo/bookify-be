package validate_data

import (
	"bookify/internal/domain"
	"bookify/pkg/shared/constants"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func ValidateEventDiscount(input domain.EventDiscount) error {
	if input.ID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}
	if input.EventID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}
	if input.DiscountValue < 0 {
		return errors.New(constants.MsgInvalidInput)
	}
	if input.DiscountUnit == "" {
		return errors.New(constants.MsgInvalidInput)
	}
	if input.DateCreated.IsZero() || input.DateCreated.Before(time.Now()) {
		return errors.New(constants.MsgInvalidInput)
	}
	if input.StartDate.IsZero() || input.StartDate.Before(time.Now()) {
		return errors.New(constants.MsgInvalidInput)
	}
	if input.EndDate.IsZero() {
		return errors.New(constants.MsgInvalidInput)
	}

	return nil
}

func ValidateEventDiscountInput(input *domain.EventDiscountInput) error {
	if input.DiscountValue < 0 {
		return errors.New(constants.MsgInvalidInput)
	}
	if input.DiscountUnit == "" {
		return errors.New(constants.MsgInvalidInput)
	}
	if input.DateCreated.IsZero() || input.DateCreated.Before(time.Now()) {
		return errors.New(constants.MsgInvalidInput)
	}
	if input.StartDate.IsZero() || input.StartDate.Before(time.Now()) {
		return errors.New(constants.MsgInvalidInput)
	}
	if input.EndDate.IsZero() {
		return errors.New(constants.MsgInvalidInput)
	}

	return nil
}
