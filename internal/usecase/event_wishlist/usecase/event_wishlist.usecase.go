package event_wishlist_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	event_repository "bookify/internal/repository/event/repository"
	eventwishlistrepository "bookify/internal/repository/event_wishlist/repository"
	userrepository "bookify/internal/repository/user/repository"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/validate_data"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IEventWishlistUseCase interface {
	GetByID(ctx context.Context, id string) (domain.EventWishlist, error)
	GetAll(ctx context.Context) ([]domain.EventWishlist, error)
	CreateOne(ctx context.Context, eventWishlist *domain.EventWishlistInput, currentUser string) error
	UpdateOne(ctx context.Context, id string, eventWishlist *domain.EventWishlistInput, currentUser string) error
	DeleteOne(ctx context.Context, id string, currentUser string) error
}

type eventWishlistUseCase struct {
	database                *config.Database
	contextTimeout          time.Duration
	eventWishlistRepository eventwishlistrepository.IEventWishlistRepository
	eventRepository         event_repository.IEventRepository
	userRepository          userrepository.IUserRepository
}

func NewEventWishlistUseCase(database *config.Database, contextTimeout time.Duration,
	eventWishlistRepository eventwishlistrepository.IEventWishlistRepository, userRepository userrepository.IUserRepository) IEventWishlistUseCase {
	return &eventWishlistUseCase{database: database, contextTimeout: contextTimeout, eventWishlistRepository: eventWishlistRepository, userRepository: userRepository}
}

func (e *eventWishlistUseCase) GetByID(ctx context.Context, id string) (domain.EventWishlist, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	eventWishlistID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.EventWishlist{}, err
	}

	data, err := e.eventWishlistRepository.GetByID(ctx, eventWishlistID)
	if err != nil {
		return domain.EventWishlist{}, err
	}

	return data, nil
}

func (e *eventWishlistUseCase) GetAll(ctx context.Context) ([]domain.EventWishlist, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.eventWishlistRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e *eventWishlistUseCase) CreateOne(ctx context.Context, eventWishlist *domain.EventWishlistInput, currentUser string) error {
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

	eventID, err := primitive.ObjectIDFromHex(currentUser)
	if err != nil {
		return err
	}

	eventData, err := e.eventRepository.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	if err = validate_data.ValidateEventWishlistInput(eventWishlist); err != nil {
		return err
	}

	eventWishlistInput := domain.EventWishlist{
		ID:        primitive.NewObjectID(),
		UserID:    userData.ID,
		EventID:   eventData.ID,
		Notes:     eventWishlist.Notes,
		CreatedAt: time.Now(),
	}

	err = e.eventWishlistRepository.CreateOne(ctx, eventWishlistInput)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventWishlistUseCase) UpdateOne(ctx context.Context, id string, eventWishlist *domain.EventWishlistInput, currentUser string) error {
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

	eventID, err := primitive.ObjectIDFromHex(currentUser)
	if err != nil {
		return err
	}

	eventData, err := e.eventRepository.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	eventWishlistID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if err = validate_data.ValidateEventWishlistInput(eventWishlist); err != nil {
		return err
	}

	eventWishlistInput := domain.EventWishlist{
		ID:      eventWishlistID,
		UserID:  userData.ID,
		EventID: eventData.ID,
		Notes:   eventWishlist.Notes,
	}

	err = e.eventWishlistRepository.UpdateOne(ctx, eventWishlistInput)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventWishlistUseCase) DeleteOne(ctx context.Context, id string, currentUser string) error {
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

	eventWishlistID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = e.eventWishlistRepository.DeleteOne(ctx, eventWishlistID)
	if err != nil {
		return err
	}

	return nil
}
