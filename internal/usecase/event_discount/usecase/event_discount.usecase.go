package event_discount_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	eventrepository "bookify/internal/repository/event/repository"
	eventdiscountrepository "bookify/internal/repository/event_discount/repository"
	userrepository "bookify/internal/repository/user/repository"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/validate_data"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IEventDiscountUseCase interface {
	GetByID(ctx context.Context, id string) (domain.EventDiscount, error)
	GetAll(ctx context.Context) ([]domain.EventDiscount, error)
	CreateOne(ctx context.Context, discount *domain.EventDiscountInput, currentUser string) error
	UpdateOne(ctx context.Context, id string, discount *domain.EventDiscountInput, currentUser string) error
	DeleteOne(ctx context.Context, id string, currentUser string) error
}

type eventDiscountUseCase struct {
	database                *config.Database
	contextTimeout          time.Duration
	eventDiscountRepository eventdiscountrepository.IEventDiscountRepository
	eventRepository         eventrepository.IEventRepository
	userRepository          userrepository.IUserRepository
}

func NewEventTypeUseCase(database *config.Database, contextTimeout time.Duration, eventDiscountRepository eventdiscountrepository.IEventDiscountRepository,
	userRepository userrepository.IUserRepository) IEventDiscountUseCase {
	return &eventDiscountUseCase{database: database, contextTimeout: contextTimeout, eventDiscountRepository: eventDiscountRepository, userRepository: userRepository}
}

func (e *eventDiscountUseCase) GetByID(ctx context.Context, id string) (domain.EventDiscount, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	eventDiscountID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.EventDiscount{}, err
	}

	data, err := e.eventDiscountRepository.GetByID(ctx, eventDiscountID)
	if err != nil {
		return domain.EventDiscount{}, err
	}

	return data, nil
}

func (e *eventDiscountUseCase) GetAll(ctx context.Context) ([]domain.EventDiscount, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.eventDiscountRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e *eventDiscountUseCase) CreateOne(ctx context.Context, discount *domain.EventDiscountInput, currentUser string) error {
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

	// check role user
	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	if err = validate_data.ValidateEventDiscountInput(discount); err != nil {
		return err
	}

	eventID, err := primitive.ObjectIDFromHex(discount.EventID)
	if err != nil {
		return err
	}

	eventData, err := e.eventRepository.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	applicableUsers, err := ConvertStringsToObjectIDs(discount.ApplicableUsers)
	if err != nil {
		return err
	}

	eventTypeInput := domain.EventDiscount{
		ID:              primitive.NewObjectID(),
		EventID:         eventData.ID,
		DiscountValue:   discount.DiscountValue,
		DiscountUnit:    discount.DiscountUnit,
		StartDate:       discount.StartDate,
		EndDate:         discount.EndDate,
		ApplicableUsers: applicableUsers,
		DateCreated:     time.Now(),
	}

	err = e.eventDiscountRepository.CreateOne(ctx, eventTypeInput)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventDiscountUseCase) UpdateOne(ctx context.Context, id string, discount *domain.EventDiscountInput, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	discountID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	userID, err := primitive.ObjectIDFromHex(currentUser)
	if err != nil {
		return err
	}

	userData, err := e.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// check role user
	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	if err = validate_data.ValidateEventDiscountInput(discount); err != nil {
		return err
	}

	eventID, err := primitive.ObjectIDFromHex(discount.EventID)
	if err != nil {
		return err
	}

	eventData, err := e.eventRepository.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	applicableUsers, err := ConvertStringsToObjectIDs(discount.ApplicableUsers)
	if err != nil {
		return err
	}

	eventTypeInput := domain.EventDiscount{
		ID:              discountID,
		EventID:         eventData.ID,
		DiscountValue:   discount.DiscountValue,
		DiscountUnit:    discount.DiscountUnit,
		StartDate:       discount.StartDate,
		EndDate:         discount.EndDate,
		ApplicableUsers: applicableUsers,
		DateCreated:     time.Now(),
	}

	err = e.eventDiscountRepository.CreateOne(ctx, eventTypeInput)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventDiscountUseCase) DeleteOne(ctx context.Context, id string, currentUser string) error {
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

	// check role user
	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	discountID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = e.eventDiscountRepository.DeleteOne(ctx, discountID)
	if err != nil {
		return err
	}

	return nil
}

func ConvertStringsToObjectIDs(ids []string) ([]primitive.ObjectID, error) {
	var objectIDs []primitive.ObjectID

	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("invalid ObjectID: %s", id)
		}
		objectIDs = append(objectIDs, objID)
	}

	return objectIDs, nil
}
