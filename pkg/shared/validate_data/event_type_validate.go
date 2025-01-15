package validate_data

import (
	"bookify/internal/domain"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateEventType(input *domain.EventType) error {
	if input.ID == primitive.NilObjectID {
		return errors.New("the event type's information cannot be empty")
	}

	if input.EventTypeName == "" {
		return errors.New("the event type's information cannot be empty")
	}

	return nil
}
