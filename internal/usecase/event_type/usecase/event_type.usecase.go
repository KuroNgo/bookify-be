package event_type_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	eventtyperepository "bookify/internal/repository/event_type/repository"
	userrepository "bookify/internal/repository/user/repository"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/validate_data"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IEventTypeUseCase interface {
	GetByID(ctx context.Context, id string) (domain.EventType, error)
	GetByName(ctx context.Context, name string) (domain.EventType, error)
	GetAll(ctx context.Context) ([]domain.EventType, error)
	CreateOne(ctx context.Context, eventType *domain.EventTypeInput, currentUser string) error
	UpdateOne(ctx context.Context, id string, eventType *domain.EventTypeInput, currentUser string) error
	DeleteOne(ctx context.Context, id string, currentUser string) error
}

type eventTypeUseCase struct {
	database            *config.Database
	contextTimeout      time.Duration
	eventTypeRepository eventtyperepository.IEventTypeRepository
	userRepository      userrepository.IUserRepository
}

func NewEventTypeUseCase(database *config.Database, contextTimeout time.Duration, eventTypeRepository eventtyperepository.IEventTypeRepository, userRepository userrepository.IUserRepository) IEventTypeUseCase {
	return &eventTypeUseCase{database: database, contextTimeout: contextTimeout, eventTypeRepository: eventTypeRepository, userRepository: userRepository}
}

func (e eventTypeUseCase) GetByName(ctx context.Context, name string) (domain.EventType, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	if name == "" {
		return domain.EventType{}, errors.New(constants.MsgInvalidInput)
	}

	data, err := e.eventTypeRepository.GetByName(ctx, name)
	if err != nil {
		return domain.EventType{}, err
	}

	return data, nil
}

func (e eventTypeUseCase) GetByID(ctx context.Context, id string) (domain.EventType, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	eventTypeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.EventType{}, err
	}

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

func (e eventTypeUseCase) CreateOne(ctx context.Context, eventType *domain.EventTypeInput, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	userID, err := primitive.ObjectIDFromHex(currentUser)
	if err != nil {
		return err
	}

	userData, err := e.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	if err = validate_data.ValidateEventTypeInput(eventType); err != nil {
		return err
	}

	count, err := e.eventTypeRepository.CountExist(ctx, eventType.Name)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(constants.MsgAPIConflict)
	}

	eventTypeInput := domain.EventType{
		ID:   primitive.NewObjectID(),
		Name: eventType.Name,
	}

	err = e.eventTypeRepository.CreateOne(ctx, eventTypeInput)
	if err != nil {
		return err
	}

	return nil
}

func (e eventTypeUseCase) UpdateOne(ctx context.Context, id string, eventType *domain.EventTypeInput, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	userID, err := primitive.ObjectIDFromHex(currentUser)
	if err != nil {
		return err
	}

	userData, err := e.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	if err := validate_data.ValidateEventTypeInput(eventType); err != nil {
		return err
	}

	count, err := e.eventTypeRepository.CountExist(ctx, eventType.Name)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(constants.MsgAPIConflict)
	}

	eventTypeId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	eventTypeInput := domain.EventType{
		ID:   eventTypeId,
		Name: eventType.Name,
	}

	err = e.eventTypeRepository.UpdateOne(ctx, eventTypeInput)
	if err != nil {
		return err
	}

	return nil
}

func (e eventTypeUseCase) DeleteOne(ctx context.Context, id string, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	userID, err := primitive.ObjectIDFromHex(currentUser)
	if err != nil {
		return err
	}

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
