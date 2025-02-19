package event_ticket_assignment_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	eventdiscountrepository "bookify/internal/repository/event_discount/repository"
	event_ticket_repository "bookify/internal/repository/event_ticket_assignment/repository"
	userrepository "bookify/internal/repository/user/repository"
	"bookify/pkg/shared/constants"
	"context"
	"errors"
	"github.com/dgraph-io/ristretto/v2"
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
	eventDiscountRepository         eventdiscountrepository.IEventDiscountRepository
	userRepository                  userrepository.IUserRepository
	cache                           *ristretto.Cache[string, domain.EventTicketAssignment]
	caches                          *ristretto.Cache[string, []domain.EventTicketAssignment]
}

// NewCacheEventTicketAssignment Kiểm tra bộ đệm khi đã đạt đến giới hạn MaxCost
// Nếu bộ nhớ vượt quá MaxCost, Ristretto sẽ tự động xóa các mục có chi phí thấp nhất
func NewCacheEventTicketAssignment() (*ristretto.Cache[string, domain.EventTicketAssignment], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, domain.EventTicketAssignment]{
		NumCounters: 1e7,       // number of keys to track frequency of (10M)
		MaxCost:     100 << 20, // 100MB // maximum cost of cache (100MB)
		BufferItems: 64,        // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

// NewCacheEventTicketAssignments Kiểm tra bộ đệm khi đã đạt đến giới hạn MaxCost
// Nếu bộ nhớ vượt quá MaxCost, Ristretto sẽ tự động xóa các mục có chi phí thấp nhất
func NewCacheEventTicketAssignments() (*ristretto.Cache[string, []domain.EventTicketAssignment], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, []domain.EventTicketAssignment]{
		NumCounters: 1e7,       // number of keys to track frequency of (10M)
		MaxCost:     100 << 20, // 100MB // maximum cost of cache (100MB)
		BufferItems: 64,        // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

func NewEventTicketAssignmentUseCase(database *config.Database, contextTimeout time.Duration,
	eventTicketAssignmentRepository event_ticket_repository.IEventTicketAssignmentRepository, userRepository userrepository.IUserRepository) IEventTicketAssignmentUseCase {
	cache, err := NewCacheEventTicketAssignment()
	if err != nil {
		panic(err)
	}

	caches, err := NewCacheEventTicketAssignments()
	if err != nil {
		panic(err)
	}
	return &eventTicketAssignmentUseCase{cache: cache, caches: caches, database: database, contextTimeout: contextTimeout,
		eventTicketAssignmentRepository: eventTicketAssignmentRepository, userRepository: userRepository}
}

func (e *eventTicketAssignmentUseCase) GetByID(ctx context.Context, id string) (domain.EventTicketAssignment, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := e.cache.Get(id)
	if found {
		return value, nil
	}

	eventTypeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.EventTicketAssignment{}, err
	}

	data, err := e.eventTicketAssignmentRepository.GetByID(ctx, eventTypeID)
	if err != nil {
		return domain.EventTicketAssignment{}, err
	}

	e.cache.Set(id, data, 1)
	e.cache.Wait()

	return data, nil
}

func (e *eventTicketAssignmentUseCase) GetAll(ctx context.Context) ([]domain.EventTicketAssignment, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := e.caches.Get("event_ticket_assignments")
	if found {
		return value, nil
	}

	data, err := e.eventTicketAssignmentRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	e.caches.Set("event_ticket_assignments", data, 1)
	e.caches.Wait()

	return data, nil
}

func (e *eventTicketAssignmentUseCase) CreateOne(ctx context.Context, eventTicketAssignment *domain.EventTicketAssignmentInput, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	userID, err := primitive.ObjectIDFromHex(currentUser)
	if err != nil {
		return err
	}

	discount, err := e.eventDiscountRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	var costsPayable float64
	if discount.DiscountUnit == "percent" {
		costsPayable = discount.DiscountValue / 100 * eventTicketAssignment.Price
	} else if discount.DiscountUnit == "amount" {
		costsPayable = eventTicketAssignment.Price - discount.DiscountValue
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
		Price:        costsPayable,
		TicketType:   eventTicketAssignment.TicketType,
		Status:       eventTicketAssignment.Status,
	}

	err = e.eventTicketAssignmentRepository.CreateOne(ctx, eventTicketAssignmentInput)
	if err != nil {
		return err
	}

	e.caches.Clear()

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

	e.caches.Clear()
	e.cache.Clear()

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

	e.caches.Clear()
	e.cache.Clear()

	return nil
}
