package event_ticket_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	eventrepository "bookify/internal/repository/event/repository"
	eventticketrepository "bookify/internal/repository/event_ticket/repository"
	userrepository "bookify/internal/repository/user/repository"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/validate_data"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IEventTicketUseCase interface {
	GetByID(ctx context.Context, id string) (domain.EventTicket, error)
	GetAll(ctx context.Context) ([]domain.EventTicket, error)
	CreateOne(ctx context.Context, eventTicket *domain.EventTicketInput, currentUser string) error
	UpdateOne(ctx context.Context, id string, eventTicket *domain.EventTicketInput, currentUser string) error
	DeleteOne(ctx context.Context, id string, currentUser string) error
}

type eventTypeUseCase struct {
	database              *config.Database
	contextTimeout        time.Duration
	eventTicketRepository eventticketrepository.IEventTicketRepository
	eventRepository       eventrepository.IEventRepository
	userRepository        userrepository.IUserRepository
}

func NewEventTicketUseCase(database *config.Database, contextTimeout time.Duration, eventTicketRepository eventticketrepository.IEventTicketRepository,
	eventRepository eventrepository.IEventRepository, userRepository userrepository.IUserRepository) IEventTicketUseCase {
	return &eventTypeUseCase{database: database, contextTimeout: contextTimeout, eventTicketRepository: eventTicketRepository, eventRepository: eventRepository, userRepository: userRepository}
}

func (e eventTypeUseCase) GetByID(ctx context.Context, id string) (domain.EventTicket, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	eventTicketID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.EventTicket{}, err
	}

	data, err := e.eventTicketRepository.GetByID(ctx, eventTicketID)
	if err != nil {
		return domain.EventTicket{}, err
	}

	return data, nil
}

func (e eventTypeUseCase) GetAll(ctx context.Context) ([]domain.EventTicket, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.eventTicketRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e eventTypeUseCase) CreateOne(ctx context.Context, eventTicket *domain.EventTicketInput, currentUser string) error {
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

	if err = validate_data.ValidateEventTicketInput(eventTicket); err != nil {
		return err
	}

	eventTicketRequest := domain.EventTicket{
		ID:       primitive.NewObjectID(),
		EventID:  eventTicket.EventID,
		Price:    eventTicket.Price,
		Quantity: eventTicket.Quantity,
		Status:   eventTicket.Status,
	}

	err = e.eventTicketRepository.CreateOne(ctx, eventTicketRequest)
	if err != nil {
		return err
	}

	return nil
}

func (e eventTypeUseCase) UpdateOne(ctx context.Context, id string, eventTicket *domain.EventTicketInput, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	eventTicketID, err := primitive.ObjectIDFromHex(id)
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

	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	if err = validate_data.ValidateEventTicketInput(eventTicket); err != nil {
		return err
	}

	eventTicketRequest := domain.EventTicket{
		ID:       eventTicketID,
		EventID:  eventTicket.EventID,
		Price:    eventTicket.Price,
		Quantity: eventTicket.Quantity,
		Status:   eventTicket.Status,
	}

	err = e.eventTicketRepository.UpdateOne(ctx, eventTicketRequest)
	if err != nil {
		return err
	}

	return nil
}

func (e eventTypeUseCase) DeleteOne(ctx context.Context, id string, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	eventTicketID, err := primitive.ObjectIDFromHex(id)
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

	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	err = e.eventTicketRepository.DeleteOne(ctx, eventTicketID)
	if err != nil {
		return err
	}

	return nil
}
