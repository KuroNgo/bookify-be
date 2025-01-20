package validate_data

import (
	"bookify/internal/domain"
	"bookify/pkg/shared/constants"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateEvent(input *domain.Event) error {
	if input.ID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.OrganizationID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.EventTypeID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.VenueID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Title == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Description == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.ImageURL == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.AssetURL == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.StartTime.IsZero() {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.EndTime.IsZero() {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Mode == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.EstimatedAttendee == 0 {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.ActualAttendee == 0 {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.TotalExpenditure == 0 {
		return errors.New(constants.MsgInvalidInput)
	}

	return nil
}

func ValidateEventInput(input *domain.EventInput) error {
	if input.OrganizationID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.EventTypeID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.VenueID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Title == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Description == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.ImageURL == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.AssetURL == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.StartTime.IsZero() {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.EndTime.IsZero() {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Mode == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.EstimatedAttendee == 0 {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.ActualAttendee == 0 {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.TotalExpenditure == 0 {
		return errors.New(constants.MsgInvalidInput)
	}

	return nil
}
