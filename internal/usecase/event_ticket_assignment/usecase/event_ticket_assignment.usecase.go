package event_ticket_assignment_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	eventrepository "bookify/internal/repository/event/repository"
	eventdiscountrepository "bookify/internal/repository/event_discount/repository"
	eventticketrepository "bookify/internal/repository/event_ticket/repository"
	eventticketassignmentrepository "bookify/internal/repository/event_ticket_assignment/repository"
	userrepository "bookify/internal/repository/user/repository"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/helper"
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
	eventTicketAssignmentRepository eventticketassignmentrepository.IEventTicketAssignmentRepository
	eventDiscountRepository         eventdiscountrepository.IEventDiscountRepository
	eventRepository                 eventrepository.IEventRepository
	userRepository                  userrepository.IUserRepository
	eventTicketRepository           eventticketrepository.IEventTicketRepository
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
	eventTicketAssignmentRepository eventticketassignmentrepository.IEventTicketAssignmentRepository,
	eventTicketRepository eventticketrepository.IEventTicketRepository, userRepository userrepository.IUserRepository) IEventTicketAssignmentUseCase {
	cache, err := NewCacheEventTicketAssignment()
	if err != nil {
		panic(err)
	}

	caches, err := NewCacheEventTicketAssignments()
	if err != nil {
		panic(err)
	}
	return &eventTicketAssignmentUseCase{cache: cache, caches: caches, database: database, contextTimeout: contextTimeout,
		eventTicketAssignmentRepository: eventTicketAssignmentRepository, eventTicketRepository: eventTicketRepository,
		userRepository: userRepository}
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

// CreateOne creates a new event ticket assignment for a user.
// If the user is eligible for a discount, the final price is adjusted accordingly.
func (e *eventTicketAssignmentUseCase) CreateOne(ctx context.Context, eventTicketAssignment *domain.EventTicketAssignmentInput, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	userID, err := primitive.ObjectIDFromHex(currentUser)
	if err != nil {
		return err
	}

	eventID, err := primitive.ObjectIDFromHex(eventTicketAssignment.EventID)
	if err != nil {
		return err
	}

	eventData, err := e.eventRepository.GetByID(ctx, eventID)
	if err != nil {
		return errors.New(constants.MsgDataNotFound)
	}

	// get discount
	discount, err := e.eventDiscountRepository.GetByUserIDInApplicableAndEventID(ctx, userID, eventID)
	if err != nil {
		return err
	}

	// calculate a value of cost payable
	var costsPayable = eventTicketAssignment.Price * float64(eventTicketAssignment.Quantity)
	if !helper.IsZeroValue(discount) { // Chỉ áp dụng giảm giá nếu có
		if discount.DiscountUnit == "percent" {
			costsPayable = eventTicketAssignment.Price * (1 - discount.DiscountValue/100)
		} else if discount.DiscountUnit == "amount" {
			costsPayable = eventTicketAssignment.Price - discount.DiscountValue
		}

		if costsPayable < 0 {
			costsPayable = 0
		}
	}

	attendanceID, err := primitive.ObjectIDFromHex(eventTicketAssignment.AttendanceID)
	if err != nil {
		return err
	}

	eventTicketAssignmentInput := domain.EventTicketAssignment{
		ID:           primitive.NewObjectID(),
		EventID:      eventData.ID,
		AttendanceID: attendanceID,
		PurchaseDate: time.Now(),
		ExpiryDate:   eventTicketAssignment.ExpiryDate,
		Price:        costsPayable,
		Quantity:     eventTicketAssignment.Quantity,
		TicketType:   eventTicketAssignment.TicketType,
		Status:       eventTicketAssignment.Status,
	}

	err = e.eventTicketAssignmentRepository.CreateOne(ctx, eventTicketAssignmentInput)
	if err != nil {
		return err
	}

	go func(eventID primitive.ObjectID) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Giới hạn 5 giây
		defer cancel()

		eventTicketData, err := e.eventTicketRepository.GetByEventID(bgCtx, eventID)
		if err != nil {
			return
		}

		_ = e.eventTicketRepository.UpdateQuantity(bgCtx, eventTicketData.ID, 1)
	}(eventID)

	// clear cache
	e.caches.Clear()

	return nil
}

// UpdateOne updates an existing event ticket assignment.
// Only a super admin can perform this action.
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
