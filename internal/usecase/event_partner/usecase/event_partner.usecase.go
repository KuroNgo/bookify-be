package event_partner_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	event_partner_repository "bookify/internal/repository/event_partner/repository"
	userrepository "bookify/internal/repository/user/repository"
	"bookify/pkg/shared/constants"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IEventPartnerUseCase interface {
	GetByID(ctx context.Context, id string) (domain.EventPartner, error)
	GetAll(ctx context.Context) ([]domain.EventPartner, error)
	CreateOne(ctx context.Context, eventPartner *domain.EventPartnerInput, currentUser string) error
	UpdateOne(ctx context.Context, id string, eventPartner *domain.EventPartnerInput, currentUser string) error
	DeleteOne(ctx context.Context, id string, currentUser string) error
}

type eventPartnerUseCase struct {
	database               *config.Database
	contextTimeout         time.Duration
	eventPartnerRepository event_partner_repository.IEventPartnerRepository
	userRepository         userrepository.IUserRepository
}

func NewEventPartnerUseCase(database *config.Database, contextTimeout time.Duration, eventPartnerRepository event_partner_repository.IEventPartnerRepository,
	userRepository userrepository.IUserRepository) IEventPartnerUseCase {
	return &eventPartnerUseCase{database: database, contextTimeout: contextTimeout,
		eventPartnerRepository: eventPartnerRepository, userRepository: userRepository}
}

func (e eventPartnerUseCase) GetByID(ctx context.Context, id string) (domain.EventPartner, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	eventPartnerID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.EventPartner{}, err
	}

	data, err := e.eventPartnerRepository.GetByID(ctx, eventPartnerID)
	if err != nil {
		return domain.EventPartner{}, err
	}

	return data, nil
}

func (e eventPartnerUseCase) GetAll(ctx context.Context) ([]domain.EventPartner, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.eventPartnerRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e eventPartnerUseCase) CreateOne(ctx context.Context, eventPartner *domain.EventPartnerInput, currentUser string) error {
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

	//if err = validate_data.ValidateEventTypeInput(eventType); err != nil {
	//	return err
	//}

	eventID, err := primitive.ObjectIDFromHex(eventPartner.EventID)
	if err != nil {
		return err
	}

	partnerID, err := primitive.ObjectIDFromHex(eventPartner.PartnerID)
	if err != nil {
		return err
	}

	eventPartnerInput := domain.EventPartner{
		ID:        primitive.NewObjectID(),
		EventID:   eventID,
		PartnerID: partnerID,
		Role:      eventPartner.Role,
	}

	err = e.eventPartnerRepository.CreateOne(ctx, eventPartnerInput)
	if err != nil {
		return err
	}

	return nil
}

func (e eventPartnerUseCase) UpdateOne(ctx context.Context, id string, eventPartner *domain.EventPartnerInput, currentUser string) error {
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

	//if err := validate_data.ValidateEventTypeInput(eventType); err != nil {
	//	return err
	//}

	eventPartnerID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	eventID, err := primitive.ObjectIDFromHex(eventPartner.EventID)
	if err != nil {
		return err
	}

	partnerID, err := primitive.ObjectIDFromHex(eventPartner.PartnerID)
	if err != nil {
		return err
	}

	eventPartnerInput := domain.EventPartner{
		ID:        eventPartnerID,
		EventID:   eventID,
		PartnerID: partnerID,
		Role:      eventPartner.Role,
	}

	err = e.eventPartnerRepository.UpdateOne(ctx, eventPartnerInput)
	if err != nil {
		return err
	}

	return nil
}

func (e eventPartnerUseCase) DeleteOne(ctx context.Context, id string, currentUser string) error {
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

	eventPartnerID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = e.eventPartnerRepository.DeleteOne(ctx, eventPartnerID)
	if err != nil {
		return err
	}

	return nil
}
