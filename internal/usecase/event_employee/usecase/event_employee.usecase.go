package event_employee_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	employeerepository "bookify/internal/repository/employee/repository"
	eventemployeerepository "bookify/internal/repository/event_employee/repository"
	userrepository "bookify/internal/repository/user/repository"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/mail/handles"
	"context"
	"errors"
	"github.com/dgraph-io/ristretto/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"sync"
	"time"
)

type IEventEmployeeUseCase interface {
	ICronjobEventEmployee // embedded interface
	GetByID(ctx context.Context, id string) (domain.EventEmployee, error)
	GetByEmployeeID(ctx context.Context, id string) (domain.EventEmployeeResponse, error)
	GetAll(ctx context.Context) ([]domain.EventEmployeeResponse, error)
	CreateOne(ctx context.Context, eventEmployee *domain.EventEmployeeInput, currentUser string) error
	UpdateOne(ctx context.Context, id string, eventEmployee *domain.EventEmployeeInput, currentUser string) error
	DeleteOne(ctx context.Context, id string, currentUser string) error
	SendQuestOfEmployeeInform(ctx context.Context) error
	DeadlineInform(ctx context.Context) error
}

type eventEmployeeUseCase struct {
	database                *config.Database
	contextTimeout          time.Duration
	eventEmployeeRepository eventemployeerepository.IEventEmployeeRepository
	employeeRepository      employeerepository.IEmployeeRepository
	userRepository          userrepository.IUserRepository
	cache                   *ristretto.Cache[string, domain.EventEmployee]
	cacheEmployee           *ristretto.Cache[string, domain.EventEmployeeResponse]
	cacheEmployees          *ristretto.Cache[string, []domain.EventEmployeeResponse]
}

// NewCacheEvent Kiểm tra bộ đệm khi đã đạt đến giới hạn MaxCost
// Nếu bộ nhớ vượt quá MaxCost, Ristretto sẽ tự động xóa các mục có chi phí thấp nhất
func NewCacheEvent() (*ristretto.Cache[string, domain.EventEmployee], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, domain.EventEmployee]{
		NumCounters: 1e7,       // number of keys to track frequency of (10M)
		MaxCost:     100 << 20, // 100MB // maximum cost of cache (100MB)
		BufferItems: 64,        // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

// NewCacheEventEmployee Kiểm tra bộ đệm khi đã đạt đến giới hạn MaxCost
// Nếu bộ nhớ vượt quá MaxCost, Ristretto sẽ tự động xóa các mục có chi phí thấp nhất
func NewCacheEventEmployee() (*ristretto.Cache[string, domain.EventEmployeeResponse], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, domain.EventEmployeeResponse]{
		NumCounters: 1e7,       // number of keys to track frequency of (10M)
		MaxCost:     100 << 20, // 100MB // maximum cost of cache (100MB)
		BufferItems: 64,        // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

// NewCacheEventEmployees Kiểm tra bộ đệm khi đã đạt đến giới hạn MaxCost
// Nếu bộ nhớ vượt quá MaxCost, Ristretto sẽ tự động xóa các mục có chi phí thấp nhất
func NewCacheEventEmployees() (*ristretto.Cache[string, []domain.EventEmployeeResponse], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, []domain.EventEmployeeResponse]{
		NumCounters: 1e7,       // number of keys to track frequency of (10M)
		MaxCost:     100 << 20, // 100MB // maximum cost of cache (100MB)
		BufferItems: 64,        // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

func NewEventEmployeeUseCase(database *config.Database, contextTimeout time.Duration, eventEmployeeRepository eventemployeerepository.IEventEmployeeRepository,
	employeeRepository employeerepository.IEmployeeRepository, userRepository userrepository.IUserRepository) IEventEmployeeUseCase {
	cache, err := NewCacheEvent()
	if err != nil {
		panic(err)
	}

	cacheEmployee, err := NewCacheEventEmployee()
	if err != nil {
		panic(err)
	}

	cacheEmployees, err := NewCacheEventEmployees()
	if err != nil {
		panic(err)
	}
	return &eventEmployeeUseCase{cache: cache, cacheEmployee: cacheEmployee, cacheEmployees: cacheEmployees, database: database, contextTimeout: contextTimeout,
		eventEmployeeRepository: eventEmployeeRepository, employeeRepository: employeeRepository, userRepository: userRepository}
}

func (e *eventEmployeeUseCase) GetByID(ctx context.Context, id string) (domain.EventEmployee, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := e.cache.Get(id)
	if found {
		return value, nil
	}

	eventTypeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.EventEmployee{}, err
	}

	data, err := e.eventEmployeeRepository.GetByID(ctx, eventTypeID)
	if err != nil {
		return domain.EventEmployee{}, err
	}

	e.cache.Set(id, data, 1)
	e.cache.Wait()

	return data, nil
}

func (e *eventEmployeeUseCase) GetByEmployeeID(ctx context.Context, id string) (domain.EventEmployeeResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := e.cacheEmployee.Get(id)
	if found {
		return value, nil
	}

	eventEmployeeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.EventEmployeeResponse{}, err
	}

	EventEmployeedata, err := e.eventEmployeeRepository.GetByID(ctx, eventEmployeeID)
	if err != nil {
		return domain.EventEmployeeResponse{}, err
	}

	uncompletedData, err := e.eventEmployeeRepository.GetIncompleteTaskPercentage(ctx, EventEmployeedata.EmployeeID)
	if err != nil {
		return domain.EventEmployeeResponse{}, err
	}

	completedData, err := e.eventEmployeeRepository.GetCompleteTaskPercentage(ctx, EventEmployeedata.EmployeeID)
	if err != nil {
		return domain.EventEmployeeResponse{}, err
	}

	response := domain.EventEmployeeResponse{
		EventEmployee:         EventEmployeedata,
		ResultEventUnComplete: uncompletedData,
		ResultEventComplete:   completedData,
	}

	e.cacheEmployee.Set(id, response, 1)
	e.cacheEmployee.Wait()

	return response, nil
}

func (e *eventEmployeeUseCase) GetAll(ctx context.Context) ([]domain.EventEmployeeResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := e.cacheEmployees.Get("event_employees")
	if found {
		return value, nil
	}

	data, err := e.eventEmployeeRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	var result []domain.EventEmployeeResponse
	result = make([]domain.EventEmployeeResponse, len(data))
	for _, i := range data {
		wg.Add(1)
		go func(i domain.EventEmployee) {
			defer wg.Done()
			EventEmployeeData, err := e.eventEmployeeRepository.GetByID(ctx, i.EmployeeID)
			if err != nil {
				errCh <- err
				return
			}

			uncompletedData, err := e.eventEmployeeRepository.GetIncompleteTaskPercentage(ctx, EventEmployeeData.EmployeeID)
			if err != nil {
				errCh <- err
				return
			}

			completedData, err := e.eventEmployeeRepository.GetCompleteTaskPercentage(ctx, EventEmployeeData.EmployeeID)
			if err != nil {
				errCh <- err
				return
			}

			response := domain.EventEmployeeResponse{
				EventEmployee:         EventEmployeeData,
				ResultEventUnComplete: uncompletedData,
				ResultEventComplete:   completedData,
			}

			result = append(result, response)
		}(i)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		e.cacheEmployees.Set("event_employees", result, 1)
		e.cacheEmployees.Wait()
		return result, nil
	}
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

	taskData := domain.Task{
		TaskName:       eventEmployee.TaskName,
		ImportantLevel: eventEmployee.ImportantLevel,
		StartDate:      eventEmployee.StartDate,
		Deadline:       eventEmployee.Deadline,
		TaskCompleted:  eventEmployee.TaskCompleted,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		WhoCreated:     userData.Email,
	}

	var tasksData []domain.Task
	tasksData = append(tasksData, taskData)

	eventEmployeeInput := &domain.EventEmployee{
		ID:         primitive.NewObjectID(),
		EventID:    eventEmployee.EventID,
		EmployeeID: eventEmployee.EmployeeID,
		Task:       tasksData,
	}

	err = e.eventEmployeeRepository.CreateOne(ctx, eventEmployeeInput)
	if err != nil {
		return err
	}

	e.cacheEmployees.Clear()

	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			return
		default:
			employeeData, err := e.employeeRepository.GetByID(ctx, eventEmployee.EmployeeID)
			if err != nil {
				return
			}

			emailData := handles.EmailData{
				EmployeeName: employeeData.FirstName + " " + employeeData.LastName,
				EventName:    "[Bookify] - Task notification",
				TaskName:     taskData.TaskName,
				Deadline:     taskData.Deadline.Format("2006-01-02"),
				AssignedBy:   userData.Email,
			}

			if err = handles.SendEmail(&emailData, employeeData.Email, "create_one_task.event_employee.html"); err != nil {
				log.Printf("Failed to send email: %v", err)
			}
		}
	}(ctx)

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
		ID:         eventEmployeeId,
		EventID:    eventEmployee.EventID,
		EmployeeID: eventEmployee.EmployeeID,
		Task:       []domain.Task{},
	}

	err = e.eventEmployeeRepository.UpdateOne(ctx, eventEmployeeInput)
	if err != nil {
		return err
	}

	e.cacheEmployees.Clear()
	e.cache.Clear()

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

	e.cacheEmployees.Clear()
	e.cache.Clear()

	return nil
}
