package event_ticket_assignment_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	event_ticket_repository "bookify/internal/repository/event_ticket_assignment/repository"
	userrepository "bookify/internal/repository/user/repository"
	"bookify/pkg/shared/constants"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IEventTicketAssignmentUseCase interface {
	GetByID(ctx context.Context, id string) (domain.EventTicketAssignment, error)
	GetAll(ctx context.Context) ([]domain.EventTicketAssignment, error)
	CreateOne(ctx context.Context, eventTicketAssignment *domain.EventTicketAssignmentInput, currentUser string) error
	UpdateOne(ctx context.Context, id string, eventTicketAssignment *domain.EventTicketAssignmentInput, currentUser string) error
	DeleteOne(ctx context.Context, id string, currentUser string) error
}

type eventTicketAssignmentUseCase struct {
	database                        *config.Database
	contextTimeout                  time.Duration
	eventTicketAssignmentRepository event_ticket_repository.IEventTicketAssignmentRepository
	userRepository                  userrepository.IUserRepository
}

func NewEventTicketAssignmentUseCase(database *config.Database, contextTimeout time.Duration,
	eventTicketAssignmentRepository event_ticket_repository.IEventTicketAssignmentRepository, userRepository userrepository.IUserRepository) IEventTicketAssignmentUseCase {
	return &eventTicketAssignmentUseCase{database: database, contextTimeout: contextTimeout,
		eventTicketAssignmentRepository: eventTicketAssignmentRepository, userRepository: userRepository}
}

func (e *eventTicketAssignmentUseCase) GetByID(ctx context.Context, id string) (domain.EventTicketAssignment, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	eventTypeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.EventTicketAssignment{}, err
	}

	data, err := e.eventTicketAssignmentRepository.GetByID(ctx, eventTypeID)
	if err != nil {
		return domain.EventTicketAssignment{}, err
	}

	return data, nil
}

func (e *eventTicketAssignmentUseCase) GetAll(ctx context.Context) ([]domain.EventTicketAssignment, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.eventTicketAssignmentRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e *eventTicketAssignmentUseCase) CreateOne(ctx context.Context, eventTicketAssignment *domain.EventTicketAssignmentInput, currentUser string) error {
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

	eventID, err := primitive.ObjectIDFromHex(eventTicketAssignment.EventID)
	if err != nil {
		return err
	}

	attendanceID, err := primitive.ObjectIDFromHex(eventTicketAssignment.AttendanceID)
	if err != nil {
		return err
	}

	eventTicketAssignmentInput := domain.EventTicketAssignment{
		ID:           primitive.NewObjectID(),
		EventID:      eventID,
		AttendanceID: attendanceID,
		PurchaseDate: time.Now(),
		ExpiryDate:   eventTicketAssignment.ExpiryDate,
		Price:        eventTicketAssignment.Price,
		TicketType:   eventTicketAssignment.TicketType,
		Status:       eventTicketAssignment.Status,
	}

	err = e.eventTicketAssignmentRepository.CreateOne(ctx, eventTicketAssignmentInput)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventTicketAssignmentUseCase) UpdateOne(ctx context.Context, id string, eventTicketAssignment *domain.EventTicketAssignmentInput, currentUser string) error {
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

	eventTicketAssignmentID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	eventID, err := primitive.ObjectIDFromHex(eventTicketAssignment.EventID)
	if err != nil {
		return err
	}

	attendanceID, err := primitive.ObjectIDFromHex(eventTicketAssignment.AttendanceID)
	if err != nil {
		return err
	}

	eventTicketAssignmentInput := domain.EventTicketAssignment{
		ID:           eventTicketAssignmentID,
		EventID:      eventID,
		AttendanceID: attendanceID,
		PurchaseDate: time.Now(),
		ExpiryDate:   eventTicketAssignment.ExpiryDate,
		Price:        eventTicketAssignment.Price,
		TicketType:   eventTicketAssignment.TicketType,
		Status:       eventTicketAssignment.Status,
	}

	err = e.eventTicketAssignmentRepository.UpdateOne(ctx, eventTicketAssignmentInput)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventTicketAssignmentUseCase) DeleteOne(ctx context.Context, id string, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	eventTicketAssignmentID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = e.eventTicketAssignmentRepository.DeleteOne(ctx, eventTicketAssignmentID)
	if err != nil {
		return err
	}

	return nil
}
