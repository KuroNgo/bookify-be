package validate_data

import (
	"bookify/internal/domain"
	"bookify/pkg/shared/constants"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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

	if input.StartTime.IsZero() {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.EndTime.IsZero() {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Mode == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	//if input.EstimatedAttendee == 0 {
	//	return errors.New(constants.MsgInvalidInput)
	//}
	//
	//if input.ActualAttendee == 0 {
	//	return errors.New(constants.MsgInvalidInput)
	//}
	//
	//if input.TotalExpenditure == 0 {
	//	return errors.New(constants.MsgInvalidInput)
	//}

	return nil
}

func ValidateEventInput(input *domain.EventInput) error {
	_, err := primitive.ObjectIDFromHex(input.OrganizationID)
	if err != nil {
		return err
	}

	if input.EventTypeName == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Title == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Description == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	parseStartTime, err := time.Parse(time.RFC3339, input.StartTime)
	if err != nil {
		return errors.New(constants.MsgInvalidInput)
	}

	parseEndTime, err := time.Parse(time.RFC3339, input.EndTime)
	if err != nil {
		return errors.New(constants.MsgInvalidInput)
	}

	if parseStartTime.IsZero() {
		return errors.New(constants.MsgInvalidInput)
	}

	if parseEndTime.IsZero() {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Mode == "" {
		return errors.New(constants.MsgInvalidInput)
	}

	if input.Capacity <= 0 {
		return errors.New(constants.MsgInvalidInput)
	}
	if input.EventMode == "Offline" {
		if input.AddressLine == "" {
			return errors.New(constants.MsgInvalidInput)
		}

		if input.City == "" {
			return errors.New(constants.MsgInvalidInput)
		}

		if input.Country == "" {
			return errors.New(constants.MsgInvalidInput)
		}
	} else if input.EventMode == "Online" {
		if input.LinkAttend == "" {
			return errors.New(constants.MsgInvalidInput)
		}

		if input.FromAttend == "" {
			return errors.New(constants.MsgInvalidInput)
		}
	}

	return nil
}
