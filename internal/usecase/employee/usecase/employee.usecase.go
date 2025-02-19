package employee_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	employeerepository "bookify/internal/repository/employee/repository"
	organizationrepository "bookify/internal/repository/organization/repository"
	userrepository "bookify/internal/repository/user/repository"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/mail/handles"
	"bookify/pkg/shared/validate_data"
	"context"
	"errors"
	"github.com/dgraph-io/ristretto/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IEmployeeUseCase interface {
	GetByID(ctx context.Context, id string) (domain.Employee, error)
	GetAll(ctx context.Context) ([]domain.Employee, error)
	CreateOne(ctx context.Context, employee *domain.EmployeeInput, currentUser string) error
	UpdateOne(ctx context.Context, id string, employee *domain.EmployeeInput, currentUser string) error
	DeleteOne(ctx context.Context, id string, currentUser string) error
	DeleteSoft(ctx context.Context, id string, currentUser string) error
	Restore(ctx context.Context, id string, currentUser string) error
}

type employeeUseCase struct {
	database               *config.Database
	contextTimeout         time.Duration
	employeeRepository     employeerepository.IEmployeeRepository
	userRepository         userrepository.IUserRepository
	organizationRepository organizationrepository.IOrganizationRepository
	cache                  *ristretto.Cache[string, domain.Employee]
	caches                 *ristretto.Cache[string, []domain.Employee]
}

// NewCacheEmployee Kiểm tra bộ đệm khi đã đạt đến giới hạn MaxCost
// Nếu bộ nhớ vượt quá MaxCost, Ristretto sẽ tự động xóa các mục có chi phí thấp nhất
func NewCacheEmployee() (*ristretto.Cache[string, domain.Employee], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, domain.Employee]{
		NumCounters: 1e7,       // number of keys to track frequency of (10M)
		MaxCost:     100 << 20, // 100MB // maximum cost of cache (100MB)
		BufferItems: 64,        // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

// NewCacheEmployees Kiểm tra bộ đệm khi đã đạt đến giới hạn MaxCost
// Nếu bộ nhớ vượt quá MaxCost, Ristretto sẽ tự động xóa các mục có chi phí thấp nhất
func NewCacheEmployees() (*ristretto.Cache[string, []domain.Employee], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, []domain.Employee]{
		NumCounters: 1e7,       // number of keys to track frequency of (10M)
		MaxCost:     100 << 20, // 100MB // maximum cost of cache (100MB)
		BufferItems: 64,        // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

func NewEmployeeUseCase(database *config.Database, contextTimeout time.Duration, employeeRepository employeerepository.IEmployeeRepository,
	userRepository userrepository.IUserRepository, organizationRepository organizationrepository.IOrganizationRepository) IEmployeeUseCase {
	cache, err := NewCacheEmployee()
	if err != nil {
		panic(err)
	}

	caches, err := NewCacheEmployees()
	if err != nil {
		panic(err)
	}
	return &employeeUseCase{cache: cache, caches: caches, database: database, contextTimeout: contextTimeout,
		employeeRepository: employeeRepository, userRepository: userRepository, organizationRepository: organizationRepository}
}

func (e employeeUseCase) GetByID(ctx context.Context, id string) (domain.Employee, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := e.cache.Get(id)
	if found {
		return value, nil
	}

	employeeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Employee{}, err
	}

	data, err := e.employeeRepository.GetByID(ctx, employeeID)
	if err != nil {
		return domain.Employee{}, err
	}

	e.cache.Set(id, data, 1)
	e.cache.Wait()

	return data, nil
}

func (e employeeUseCase) GetAll(ctx context.Context) ([]domain.Employee, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := e.caches.Get("employees")
	if found {
		return value, nil
	}

	data, err := e.employeeRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	e.caches.Set("employees", data, 1)
	e.caches.Wait()

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

	if err = validate_data.ValidateEmployeeInput(employee); err != nil {
		return err
	}

	organizationData, err := e.organizationRepository.GetByID(ctx, employee.OrganizationID)
	if err != nil {
		return err
	}

	employeeInput := &domain.Employee{
		ID:             primitive.NewObjectID(),
		OrganizationID: employee.OrganizationID,
		FirstName:      employee.FirstName,
		LastName:       employee.LastName,
		JobTitle:       employee.JobTitle,
		Email:          employee.Email,
		Status:         "enabled",
		UpdatedAt:      time.Now(),
		CreatedAt:      time.Now(),
		WhoCreated:     userData.Email,
	}

	err = e.employeeRepository.CreateOne(ctx, employeeInput)
	if err != nil {
		return err
	}

	e.caches.Clear()

	time.AfterFunc(time.Minute*5, func() {
		emailData := handles.EmailData{
			FullName:     employee.FirstName + " " + employee.LastName,
			Subject:      "[Bookify] - Welcome to Bookify! Your employee account is ready",
			JobTitle:     employee.JobTitle,
			Organization: organizationData.Name,
		}

		err = handles.SendEmail(&emailData, employee.Email, "create_one.employee.html")
		if err != nil {
			return
		}
	})

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

	if err = validate_data.ValidateEmployeeInput(employee); err != nil {
		return err
	}

	employeeInput := &domain.Employee{
		ID:             primitive.NewObjectID(),
		OrganizationID: employee.OrganizationID,
		FirstName:      employee.FirstName,
		LastName:       employee.LastName,
		JobTitle:       employee.JobTitle,
		Email:          employee.Email,
		UpdatedAt:      time.Now(),
	}

	err = e.employeeRepository.UpdateOne(ctx, employeeInput)
	if err != nil {
		return err
	}

	e.caches.Clear()
	e.cache.Clear()

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

	e.caches.Clear()
	e.cache.Clear()

	return nil
}

func (e employeeUseCase) DeleteSoft(ctx context.Context, id string, currentUser string) error {
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

	err = e.employeeRepository.DeleteSoft(ctx, employeeID)
	if err != nil {
		return err
	}

	e.caches.Clear()
	e.cache.Clear()

	time.AfterFunc(time.Minute*5, func() {
		employeeData, err := e.employeeRepository.GetByID(ctx, employeeID)
		if err != nil {
			return
		}

		organizationData, err := e.organizationRepository.GetByID(ctx, employeeData.OrganizationID)
		if err != nil {
			return
		}

		emailData := handles.EmailData{
			FullName:     employeeData.FirstName + " " + employeeData.LastName,
			Subject:      "[Bookify] - Employee account removal notification",
			JobTitle:     employeeData.JobTitle,
			Organization: organizationData.Name,
		}

		err = handles.SendEmail(&emailData, employeeData.Email, "delete_one.employee.html")
		if err != nil {
			return
		}
	})

	return nil
}

func (e employeeUseCase) Restore(ctx context.Context, id string, currentUser string) error {
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

	err = e.employeeRepository.Restore(ctx, employeeID)
	if err != nil {
		return err
	}

	e.caches.Clear()
	e.cache.Clear()

	// AfterFunc will call goroutine for active function with time settings
	time.AfterFunc(time.Minute*5, func() {
		employeeData, err := e.employeeRepository.GetByID(ctx, employeeID)
		if err != nil {
			return
		}

		organizationData, err := e.organizationRepository.GetByID(ctx, employeeData.OrganizationID)
		if err != nil {
			return
		}

		emailData := handles.EmailData{
			FullName:     employeeData.FirstName + " " + employeeData.LastName,
			Subject:      "[Bookify] - Employee account removal notification",
			JobTitle:     employeeData.JobTitle,
			Organization: organizationData.Name,
		}

		err = handles.SendEmail(&emailData, employeeData.Email, "restore_one.employee.html")
		if err != nil {
			return
		}
	})

	return nil
}
