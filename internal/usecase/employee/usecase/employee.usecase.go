package employee_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	employee_repository "bookify/internal/repository/employee/repository"
	userrepository "bookify/internal/repository/user/repository"
	"bookify/pkg/shared/constants"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IEmployeeUseCase interface {
	GetByID(ctx context.Context, id string) (domain.Employee, error)
	GetAll(ctx context.Context) ([]domain.Employee, error)
	CreateOne(ctx context.Context, employee *domain.EmployeeInput, currentUser string) error
	UpdateOne(ctx context.Context, id string, employee *domain.EmployeeInput, currentUser string) error
	DeleteOne(ctx context.Context, id string, currentUser string) error
}

type employeeUseCase struct {
	database           *config.Database
	contextTimeout     time.Duration
	employeeRepository employee_repository.IEmployeeRepository
	userRepository     userrepository.IUserRepository
}

func (e employeeUseCase) GetByID(ctx context.Context, id string) (domain.Employee, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	employeeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Employee{}, err
	}

	data, err := e.employeeRepository.GetByID(ctx, employeeID)
	if err != nil {
		return domain.Employee{}, err
	}

	return data, nil
}

func (e employeeUseCase) GetAll(ctx context.Context) ([]domain.Employee, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.employeeRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e employeeUseCase) CreateOne(ctx context.Context, employee *domain.EmployeeInput, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// Convert currentUser from string to primitive objectID
	userID, err := primitive.ObjectIDFromHex(currentUser)
	if err != nil {
		return err
	}

	// Handle get by id to get user data
	userData, err := e.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// from user data, check role of user
	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	//if err = validate_data.ValidatePartnerInput(venue); err != nil {
	//	return err
	//}
	//
	//count, err := p.partnerRepository.CountExist(ctx, partner.Name)
	//if err != nil {
	//	return err
	//}
	//
	//if count > 0 {
	//	return errors.New(constants.MsgAPIConflict)
	//}

	employeeInput := &domain.Employee{
		ID:             primitive.NewObjectID(),
		OrganizationID: primitive.NewObjectID(),
		FirstName:      employee.FirstName,
		LastName:       employee.LastName,
		JobTitle:       employee.JobTitle,
		Email:          employee.Email,
	}

	err = e.employeeRepository.CreateOne(ctx, employeeInput)
	if err != nil {
		return err
	}

	return nil
}

func (e employeeUseCase) UpdateOne(ctx context.Context, id string, employee *domain.EmployeeInput, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// Convert currentUser from string to primitive objectID
	userID, err := primitive.ObjectIDFromHex(currentUser)
	if err != nil {
		return err
	}

	// Handle get by id to get user data
	userData, err := e.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// from user data, check role of user
	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	//if err = validate_data.ValidatePartnerInput(venue); err != nil {
	//	return err
	//}
	//
	//count, err := p.partnerRepository.CountExist(ctx, partner.Name)
	//if err != nil {
	//	return err
	//}
	//
	//if count > 0 {
	//	return errors.New(constants.MsgAPIConflict)
	//}

	employeeInput := &domain.Employee{
		ID:             primitive.NewObjectID(),
		OrganizationID: primitive.NewObjectID(),
		FirstName:      employee.FirstName,
		LastName:       employee.LastName,
		JobTitle:       employee.JobTitle,
		Email:          employee.Email,
	}

	err = e.employeeRepository.CreateOne(ctx, employeeInput)
	if err != nil {
		return err
	}

	return nil
}

func (e employeeUseCase) DeleteOne(ctx context.Context, id string, currentUser string) error {
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

	employeeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = e.employeeRepository.DeleteOne(ctx, employeeID)
	if err != nil {
		return err
	}

	return nil
}

func NewEmployeeUseCase(database *config.Database, contextTimeout time.Duration, employeeRepository employee_repository.IEmployeeRepository,
	userRepository userrepository.IUserRepository) IEmployeeUseCase {
	return &employeeUseCase{database: database, contextTimeout: contextTimeout, employeeRepository: employeeRepository, userRepository: userRepository}
}
