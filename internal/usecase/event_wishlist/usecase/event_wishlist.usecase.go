package event_wishlist_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	event_repository "bookify/internal/repository/event/repository"
	eventwishlistrepository "bookify/internal/repository/event_wishlist/repository"
	userrepository "bookify/internal/repository/user/repository"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/validate_data"
	"context"
	"errors"
	"github.com/dgraph-io/ristretto/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"time"
)

type IEventWishlistUseCase interface {
	GetByID(ctx context.Context, id string) (domain.EventWishlist, error)
	GetAll(ctx context.Context) ([]domain.EventWishlist, error)
	CreateOne(ctx context.Context, eventWishlist *domain.EventWishlistInput, currentUser string) error
	UpdateOne(ctx context.Context, id string, eventWishlist *domain.EventWishlistInput, currentUser string) error
	DeleteOne(ctx context.Context, id string, currentUser string) error
}

type eventWishlistUseCase struct {
	database                *config.Database
	contextTimeout          time.Duration
	eventWishlistRepository eventwishlistrepository.IEventWishlistRepository
	eventRepository         event_repository.IEventRepository
	userRepository          userrepository.IUserRepository
	mu                      *sync.Mutex
	cache                   *ristretto.Cache[string, domain.EventWishlist]
	caches                  *ristretto.Cache[string, []domain.EventWishlist]
}

// NewCacheEventWishlist Kiểm tra bộ đệm khi đã đạt đến giới hạn MaxCost
// Nếu bộ nhớ vượt quá MaxCost, Ristretto sẽ tự động xóa các mục có chi phí thấp nhất
func NewCacheEventWishlist() (*ristretto.Cache[string, domain.EventWishlist], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, domain.EventWishlist]{
		NumCounters: 1e7,       // number of keys to track frequency of (10M)
		MaxCost:     100 << 20, // 100MB // maximum cost of cache (100MB)
		BufferItems: 64,        // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

// NewCacheEventWishlists Kiểm tra bộ đệm khi đã đạt đến giới hạn MaxCost
// Nếu bộ nhớ vượt quá MaxCost, Ristretto sẽ tự động xóa các mục có chi phí thấp nhất
func NewCacheEventWishlists() (*ristretto.Cache[string, []domain.EventWishlist], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, []domain.EventWishlist]{
		NumCounters: 1e7,       // number of keys to track frequency of (10M)
		MaxCost:     100 << 20, // 100MB // maximum cost of cache (100MB)
		BufferItems: 64,        // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

func NewEventWishlistUseCase(database *config.Database, contextTimeout time.Duration,
	eventWishlistRepository eventwishlistrepository.IEventWishlistRepository, userRepository userrepository.IUserRepository) IEventWishlistUseCase {
	cache, err := NewCacheEventWishlist()
	if err != nil {
		panic(err)
	}

	caches, err := NewCacheEventWishlists()
	if err != nil {
		panic(err)
	}
	return &eventWishlistUseCase{cache: cache, caches: caches, database: database, contextTimeout: contextTimeout, eventWishlistRepository: eventWishlistRepository, userRepository: userRepository}
}

func (e *eventWishlistUseCase) GetByID(ctx context.Context, id string) (domain.EventWishlist, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := e.cache.Get(id)
	if found {
		return value, nil
	}

	eventWishlistID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.EventWishlist{}, err
	}

	data, err := e.eventWishlistRepository.GetByID(ctx, eventWishlistID)
	if err != nil {
		return domain.EventWishlist{}, err
	}

	e.cache.Set(id, data, 1)
	// wait for value to pass through buffers
	e.cache.Wait()

	return data, nil
}

func (e *eventWishlistUseCase) GetAll(ctx context.Context) ([]domain.EventWishlist, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := e.caches.Get("wishlists")
	if found {
		return value, nil
	}

	data, err := e.eventWishlistRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	e.caches.Set("wishlists", data, 1)
	// wait for value to pass through buffers
	e.caches.Wait()

	return data, nil
}

func (e *eventWishlistUseCase) CreateOne(ctx context.Context, eventWishlist *domain.EventWishlistInput, currentUser string) error {
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

	eventID, err := primitive.ObjectIDFromHex(currentUser)
	if err != nil {
		return err
	}

	eventData, err := e.eventRepository.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	if err = validate_data.ValidateEventWishlistInput(eventWishlist); err != nil {
		return err
	}

	eventWishlistInput := domain.EventWishlist{
		ID:        primitive.NewObjectID(),
		UserID:    userData.ID,
		EventID:   eventData.ID,
		Notes:     eventWishlist.Notes,
		CreatedAt: time.Now(),
	}

	err = e.eventWishlistRepository.CreateOne(ctx, eventWishlistInput)
	if err != nil {
		return err
	}

	e.caches.Clear()

	return nil
}

func (e *eventWishlistUseCase) UpdateOne(ctx context.Context, id string, eventWishlist *domain.EventWishlistInput, currentUser string) error {
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

	eventID, err := primitive.ObjectIDFromHex(currentUser)
	if err != nil {
		return err
	}

	eventData, err := e.eventRepository.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	eventWishlistID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if err = validate_data.ValidateEventWishlistInput(eventWishlist); err != nil {
		return err
	}

	eventWishlistInput := domain.EventWishlist{
		ID:      eventWishlistID,
		UserID:  userData.ID,
		EventID: eventData.ID,
		Notes:   eventWishlist.Notes,
	}

	err = e.eventWishlistRepository.UpdateOne(ctx, eventWishlistInput)
	if err != nil {
		return err
	}

	e.caches.Clear()
	e.cache.Clear()

	return nil
}

func (e *eventWishlistUseCase) DeleteOne(ctx context.Context, id string, currentUser string) error {
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

	eventWishlistID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = e.eventWishlistRepository.DeleteOne(ctx, eventWishlistID)
	if err != nil {
		return err
	}

	e.caches.Clear()
	e.cache.Clear()

	return nil
}
