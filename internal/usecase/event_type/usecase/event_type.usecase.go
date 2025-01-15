package event_type_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	event_type_repository "bookify/internal/repository/event_type/repository"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/validate_data"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IEventTypeRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (domain.EventType, error)
	GetAll(ctx context.Context) ([]domain.EventType, error)
	CreateOne(ctx context.Context, eventType *domain.EventType) error
	UpdateOne(ctx context.Context, eventType *domain.EventType) error
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
}

type eventTypeUseCase struct {
	database            *config.Database
	contextTimeout      time.Duration
	eventTypeRepository event_type_repository.IEventTypeRepository
}

func (e eventTypeUseCase) GetByID(ctx context.Context, id primitive.ObjectID) (domain.EventType, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.eventTypeRepository.GetByID(ctx, id)
	if err != nil {
		return domain.EventType{}, err
	}

	return data, nil
}

func (e eventTypeUseCase) GetAll(ctx context.Context) ([]domain.EventType, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.eventTypeRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e eventTypeUseCase) CreateOne(ctx context.Context, eventType *domain.EventType) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	if err := validate_data.ValidateEventType(eventType); err != nil {
		return err
	}

	count, err := e.eventTypeRepository.CountExist(ctx, eventType.EventTypeName)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(constants.MsgAPIConflict)
	}

	err = e.eventTypeRepository.CreateOne(ctx, eventType)
	if err != nil {
		return err
	}

	return nil
}

func (e eventTypeUseCase) UpdateOne(ctx context.Context, eventType *domain.EventType) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	if err := validate_data.ValidateEventType(eventType); err != nil {
		return err
	}

	count, err := e.eventTypeRepository.CountExist(ctx, eventType.EventTypeName)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(constants.MsgAPIConflict)
	}

	err = e.eventTypeRepository.UpdateOne(ctx, eventType)
	if err != nil {
		return err
	}

	return nil
}

func (e eventTypeUseCase) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.eventTypeRepository.DeleteOne(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func NewEventTypeUseCase(database *config.Database, contextTimeout time.Duration, eventTypeRepository event_type_repository.IEventTypeRepository) IEventTypeRepository {
	return &eventTypeUseCase{database: database, contextTimeout: contextTimeout, eventTypeRepository: eventTypeRepository}
}
