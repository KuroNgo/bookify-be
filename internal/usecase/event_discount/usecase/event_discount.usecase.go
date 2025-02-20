package event_discount_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	eventrepository "bookify/internal/repository/event/repository"
	eventdiscountrepository "bookify/internal/repository/event_discount/repository"
	eventticketrepository "bookify/internal/repository/event_ticket/repository"
	eventwishlistrepository "bookify/internal/repository/event_wishlist/repository"
	userrepository "bookify/internal/repository/user/repository"
	"bookify/pkg/shared/constants"
	cronjob "bookify/pkg/shared/cron"
	"bookify/pkg/shared/validate_data"
	"context"
	"errors"
	"fmt"
	"github.com/dgraph-io/ristretto/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

type IEventDiscountUseCase interface {
	IJobWorkerEventDiscount
	GetByID(ctx context.Context, id string) (domain.EventDiscount, error)
	GetAll(ctx context.Context) ([]domain.EventDiscount, error)
	CreateOne(ctx context.Context, discount *domain.EventDiscountInput, currentUser string) error
	UpdateOne(ctx context.Context, id string, discount *domain.EventDiscountInput, currentUser string) error
	DeleteOne(ctx context.Context, id string, currentUser string) error
}

type eventDiscountUseCase struct {
	database                *config.Database
	contextTimeout          time.Duration
	cs                      *cronjob.CronScheduler
	eventDiscountRepository eventdiscountrepository.IEventDiscountRepository
	eventRepository         eventrepository.IEventRepository
	eventTicketRepository   eventticketrepository.IEventTicketRepository
	wishlistRepository      eventwishlistrepository.IEventWishlistRepository
	userRepository          userrepository.IUserRepository
	cache                   *ristretto.Cache[string, domain.EventDiscount]
	caches                  *ristretto.Cache[string, []domain.EventDiscount]
}

// NewCacheEventDiscount Kiểm tra bộ đệm khi đã đạt đến giới hạn MaxCost
// Nếu bộ nhớ vượt quá MaxCost, Ristretto sẽ tự động xóa các mục có chi phí thấp nhất
func NewCacheEventDiscount() (*ristretto.Cache[string, domain.EventDiscount], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, domain.EventDiscount]{
		NumCounters: 1e7,       // number of keys to track frequency of (10M)
		MaxCost:     100 << 20, // 100MB // maximum cost of cache (100MB)
		BufferItems: 64,        // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

// NewCacheEventDiscounts Kiểm tra bộ đệm khi đã đạt đến giới hạn MaxCost
// Nếu bộ nhớ vượt quá MaxCost, Ristretto sẽ tự động xóa các mục có chi phí thấp nhất
func NewCacheEventDiscounts() (*ristretto.Cache[string, []domain.EventDiscount], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, []domain.EventDiscount]{
		NumCounters: 1e7,       // number of keys to track frequency of (10M)
		MaxCost:     100 << 20, // 100MB // maximum cost of cache (100MB)
		BufferItems: 64,        // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

func NewEventDiscountUseCase(database *config.Database, cs *cronjob.CronScheduler, contextTimeout time.Duration, eventDiscountRepository eventdiscountrepository.IEventDiscountRepository,
	userRepository userrepository.IUserRepository) IEventDiscountUseCase {
	cache, err := NewCacheEventDiscount()
	if err != nil {
		panic(err)
	}

	caches, err := NewCacheEventDiscounts()
	if err != nil {
		panic(err)
	}
	return &eventDiscountUseCase{cache: cache, caches: caches, database: database, contextTimeout: contextTimeout,
		eventDiscountRepository: eventDiscountRepository, userRepository: userRepository, cs: cs}
}

func (e *eventDiscountUseCase) GetByID(ctx context.Context, id string) (domain.EventDiscount, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := e.cache.Get(id)
	if found {
		return value, nil
	}

	eventDiscountID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.EventDiscount{}, err
	}

	data, err := e.eventDiscountRepository.GetByID(ctx, eventDiscountID)
	if err != nil {
		return domain.EventDiscount{}, err
	}

	e.cache.Set(id, data, 1)
	e.cache.Wait()

	return data, nil
}

func (e *eventDiscountUseCase) GetAll(ctx context.Context) ([]domain.EventDiscount, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := e.caches.Get("discounts")
	if found {
		return value, nil
	}

	data, err := e.eventDiscountRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	e.caches.Set("discounts", data, 1)
	e.caches.Wait()

	return data, nil
}

// CreateOne
// - The discount's StartDate and EndDate must fall within the event's StartDate and EndDate.
// - If the discount unit is "amount", the DiscountValue must not exceed the event ticket price.
// - If the discount unit is "percent", the DiscountValue must be between 0 and 100.
func (e *eventDiscountUseCase) CreateOne(ctx context.Context, discount *domain.EventDiscountInput, currentUser string) error {
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

	// check role user
	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	if err = validate_data.ValidateEventDiscountInput(discount); err != nil {
		return err
	}

	eventID, err := primitive.ObjectIDFromHex(discount.EventID)
	if err != nil {
		return err
	}

	eventData, err := e.eventRepository.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	eventTicketData, err := e.eventTicketRepository.GetByID(ctx, eventData.ID)
	if err != nil {
		return err
	}

	// you can not create a discount with price greater than price of event
	if discount.DiscountUnit == "amount" && discount.DiscountValue > eventTicketData.Price {
		return errors.New(constants.MsgInvalidInput)
	}

	if discount.DiscountUnit == "percent" && (discount.DiscountValue > 100 || discount.DiscountValue < 0) {
		return errors.New(constants.MsgInvalidInput)
	}

	applicableUsers, err := ConvertStringsToObjectIDs(discount.ApplicableUsers)
	if err != nil {
		return err
	}

	eventTypeInput := domain.EventDiscount{
		ID:              primitive.NewObjectID(),
		EventID:         eventData.ID,
		DiscountValue:   discount.DiscountValue,
		DiscountUnit:    discount.DiscountUnit,
		StartDate:       discount.StartDate,
		EndDate:         discount.EndDate,
		ApplicableUsers: applicableUsers,
		DateCreated:     time.Now(),
	}

	err = e.eventDiscountRepository.CreateOne(ctx, eventTypeInput)
	if err != nil {
		return err
	}

	// background job
	go func() {
		err = e.JobWorkerSendDiscountForApplicableUsersIfTheyHaveWishlist(e.cs)
		if err != nil {
			log.Println("Failed to execute job:", err)
		}
	}()

	go func() {
		err = e.JobWorkerDiscountForApplicableUsersIfTheyHaveWishlistExpiringOneDayLeft(e.cs)
		if err != nil {
			log.Println("Failed to execute job:", err)
		}
	}()

	e.caches.Clear()

	return nil
}

// UpdateOne
// - Ensure the discount's StartDate and EndDate are within the event's StartDate and EndDate.
// - If the discount unit is "amount", the DiscountValue must not exceed the event ticket price.
// - If the discount unit is "percent", the DiscountValue must be between 0 and 100.
// - If the discount's StartDate changes, remove the scheduled background jobs and reinitialize them.
func (e *eventDiscountUseCase) UpdateOne(ctx context.Context, id string, discount *domain.EventDiscountInput, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	discountID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	discountData, err := e.eventDiscountRepository.GetByID(ctx, discountID)
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

	// check role user
	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	if err = validate_data.ValidateEventDiscountInput(discount); err != nil {
		return err
	}

	eventID, err := primitive.ObjectIDFromHex(discount.EventID)
	if err != nil {
		return err
	}

	eventData, err := e.eventRepository.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	eventTicketData, err := e.eventTicketRepository.GetByID(ctx, eventData.ID)
	if err != nil {
		return err
	}

	// you can not create a discount with price greater than price of event\
	if discount.DiscountUnit == "amount" && discount.DiscountValue > eventTicketData.Price {
		return errors.New(constants.MsgInvalidInput)
	}

	if discount.DiscountUnit == "percent" && (discount.DiscountValue > 100 || discount.DiscountValue < 0) {
		return errors.New(constants.MsgInvalidInput)
	}

	applicableUsers, err := ConvertStringsToObjectIDs(discount.ApplicableUsers)
	if err != nil {
		return err
	}

	eventDiscountInput := domain.EventDiscount{
		ID:              discountID,
		EventID:         eventData.ID,
		DiscountValue:   discount.DiscountValue,
		DiscountUnit:    discount.DiscountUnit,
		StartDate:       discount.StartDate,
		EndDate:         discount.EndDate,
		ApplicableUsers: applicableUsers,
		DateCreated:     time.Now(),
	}

	if !discount.StartDate.Equal(discountData.StartDate) {
		go func() {
			err = e.RemoveJobWorkerSendDiscountForApplicableUsersIfTheyHaveWishlist(e.cs)
			if err != nil {
				log.Printf("Failed to remove job for discount ID %s: %v", id, err)
			}
		}()

		go func() {
			err = e.RemoveJobWorkerDiscountForApplicableUsersIfTheyHaveWishlistExpiringOneDayLeft(e.cs)
			if err != nil {
				log.Printf("Failed to remove job for discount ID %s: %v", id, err)
			}
		}()
	}

	err = e.eventDiscountRepository.UpdateOne(ctx, eventDiscountInput)
	if err != nil {
		return err
	}

	// background job
	go func() {
		err = e.JobWorkerSendDiscountForApplicableUsersIfTheyHaveWishlist(e.cs)
		if err != nil {
			log.Println("Failed to execute job:", err)
		}
	}()

	go func() {
		err = e.JobWorkerDiscountForApplicableUsersIfTheyHaveWishlistExpiringOneDayLeft(e.cs)
		if err != nil {
			log.Println("Failed to execute job:", err)
		}
	}()

	e.caches.Clear()
	e.cache.Clear()

	return nil
}

func (e *eventDiscountUseCase) DeleteOne(ctx context.Context, id string, currentUser string) error {
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

	// check role user
	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	discountID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = e.eventDiscountRepository.DeleteOne(ctx, discountID)
	if err != nil {
		return err
	}

	go func() {
		err = e.RemoveJobWorkerSendDiscountForApplicableUsersIfTheyHaveWishlist(e.cs)
		if err != nil {
			log.Println("Failed to execute job:", err)
		}
	}()

	go func() {
		err = e.RemoveJobWorkerDiscountForApplicableUsersIfTheyHaveWishlistExpiringOneDayLeft(e.cs)
		if err != nil {
			log.Println("Failed to execute job:", err)
		}
	}()

	e.caches.Clear()
	e.cache.Clear()

	return nil
}

func ConvertStringsToObjectIDs(ids []string) ([]primitive.ObjectID, error) {
	var objectIDs []primitive.ObjectID

	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("invalid ObjectID: %s", id)
		}
		objectIDs = append(objectIDs, objID)
	}

	return objectIDs, nil
}
