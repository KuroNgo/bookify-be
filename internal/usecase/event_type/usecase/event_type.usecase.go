package event_type_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	event_type_repository "bookify/internal/repository/event_type/repository"
	user_repository "bookify/internal/repository/user/repository"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/validate_data"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IEventTypeRepository interface {
	GetByID(ctx context.Context, id string) (domain.EventType, error)
	GetAll(ctx context.Context) ([]domain.EventType, error)
	CreateOne(ctx context.Context, eventType *domain.EventType, currentUser string) error
	UpdateOne(ctx context.Context, eventType *domain.EventType, currentUser string) error
	DeleteOne(ctx context.Context, id string, currentUser string) error
}

type eventTypeUseCase struct {
	database            *config.Database
	contextTimeout      time.Duration
	eventTypeRepository event_type_repository.IEventTypeRepository
	userRepository      user_repository.IUserRepository
}

func (e eventTypeUseCase) GetByID(ctx context.Context, id string) (domain.EventType, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	eventTypeID, err := primitive.ObjectIDFromHex(id)
	data, err := e.eventTypeRepository.GetByID(ctx, eventTypeID)
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

func (e eventTypeUseCase) CreateOne(ctx context.Context, eventType *domain.EventType, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	userID, err := primitive.ObjectIDFromHex(currentUser)
	userData, err := e.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	if err = validate_data.ValidateEventType(eventType); err != nil {
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

func (e eventTypeUseCase) UpdateOne(ctx context.Context, eventType *domain.EventType, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	userID, err := primitive.ObjectIDFromHex(currentUser)
	userData, err := e.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

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

func (e eventTypeUseCase) DeleteOne(ctx context.Context, id string, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	userID, err := primitive.ObjectIDFromHex(currentUser)
	userData, err := e.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	eventTypeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = e.eventTypeRepository.DeleteOne(ctx, eventTypeID)
	if err != nil {
		return err
	}

	return nil
}

func NewEventTypeUseCase(database *config.Database, contextTimeout time.Duration, eventTypeRepository event_type_repository.IEventTypeRepository) IEventTypeRepository {
	return &eventTypeUseCase{database: database, contextTimeout: contextTimeout, eventTypeRepository: eventTypeRepository}
}
