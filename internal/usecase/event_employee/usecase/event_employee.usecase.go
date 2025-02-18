package event_employee_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	event_employee_repository "bookify/internal/repository/event_employee/repository"
	userrepository "bookify/internal/repository/user/repository"
	"bookify/pkg/shared/constants"
	"context"
	"errors"
	"github.com/dgraph-io/ristretto/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IEventEmployeeUseCase interface {
	ICronjobEventEmployee // embedded interface
	GetByID(ctx context.Context, id string) (domain.EventEmployee, error)
	GetAll(ctx context.Context) ([]domain.EventEmployee, error)
	CreateOne(ctx context.Context, eventEmployee *domain.EventEmployeeInput, currentUser string) error
	UpdateOne(ctx context.Context, id string, eventEmployee *domain.EventEmployeeInput, currentUser string) error
	DeleteOne(ctx context.Context, id string, currentUser string) error
	SendQuestOfEmployeeInform(ctx context.Context) error
	DeadlineInform(ctx context.Context) error
}

type eventEmployeeUseCase struct {
	database                *config.Database
	contextTimeout          time.Duration
	eventEmployeeRepository event_employee_repository.IEventEmployeeRepository
	userRepository          userrepository.IUserRepository
	cache                   *ristretto.Cache[string, domain.EventType]
	caches                  *ristretto.Cache[string, []domain.EventType]
}

func NewEventEmployeeUseCase(database *config.Database, contextTimeout time.Duration, eventEmployeeRepository event_employee_repository.IEventEmployeeRepository,
	userRepository userrepository.IUserRepository) IEventEmployeeUseCase {
	return &eventEmployeeUseCase{database: database, contextTimeout: contextTimeout, eventEmployeeRepository: eventEmployeeRepository, userRepository: userRepository}
}

func (e *eventEmployeeUseCase) GetByID(ctx context.Context, id string) (domain.EventEmployee, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	eventTypeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.EventEmployee{}, err
	}

	data, err := e.eventEmployeeRepository.GetByID(ctx, eventTypeID)
	if err != nil {
		return domain.EventEmployee{}, err
	}

	return data, nil
}

func (e *eventEmployeeUseCase) GetAll(ctx context.Context) ([]domain.EventEmployee, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.eventEmployeeRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e *eventEmployeeUseCase) CreateOne(ctx context.Context, eventEmployee *domain.EventEmployeeInput, currentUser string) error {
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

	eventEmployeeInput := &domain.EventEmployee{
		ID:            primitive.NewObjectID(),
		EventID:       eventEmployee.EventID,
		EmployeeID:    eventEmployee.EmployeeID,
		Task:          eventEmployee.Task,
		StartDate:     eventEmployee.Deadline,
		Deadline:      eventEmployee.Deadline,
		TaskCompleted: eventEmployee.TaskCompleted,
	}

	err = e.eventEmployeeRepository.CreateOne(ctx, eventEmployeeInput)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventEmployeeUseCase) UpdateOne(ctx context.Context, id string, eventEmployee *domain.EventEmployeeInput, currentUser string) error {
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

	eventEmployeeId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	eventEmployeeInput := &domain.EventEmployee{
		ID:            eventEmployeeId,
		EventID:       eventEmployee.EventID,
		EmployeeID:    eventEmployee.EmployeeID,
		Task:          eventEmployee.Task,
		StartDate:     eventEmployee.Deadline,
		Deadline:      eventEmployee.Deadline,
		TaskCompleted: eventEmployee.TaskCompleted,
	}

	err = e.eventEmployeeRepository.UpdateOne(ctx, eventEmployeeInput)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventEmployeeUseCase) DeleteOne(ctx context.Context, id string, currentUser string) error {
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

	eventEmployeeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = e.eventEmployeeRepository.DeleteOne(ctx, eventEmployeeID)
	if err != nil {
		return err
	}

	return nil
}
