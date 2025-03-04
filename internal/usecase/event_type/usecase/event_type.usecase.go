package event_type_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	eventtyperepository "bookify/internal/repository/event_type/repository"
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

type IEventTypeUseCase interface {
	GetByID(ctx context.Context, id string) (domain.EventType, error)
	GetByName(ctx context.Context, name string) (domain.EventType, error)
	GetAll(ctx context.Context) ([]domain.EventType, error)
	CreateOne(ctx context.Context, eventType *domain.EventTypeInput, currentUser string) error
	UpdateOne(ctx context.Context, id string, eventType *domain.EventTypeInput, currentUser string) error
	DeleteOne(ctx context.Context, id string, currentUser string) error
}

type eventTypeUseCase struct {
	database            *config.Database
	contextTimeout      time.Duration
	eventTypeRepository eventtyperepository.IEventTypeRepository
	userRepository      userrepository.IUserRepository
	mu                  *sync.Mutex
	cache               *ristretto.Cache[string, domain.EventType]
	caches              *ristretto.Cache[string, []domain.EventType]
}

// NewCache Kiểm tra bộ đệm khi đã đạt đến giới hạn MaxCost
// Nếu bộ nhớ vượt quá MaxCost, Ristretto sẽ tự động xóa các mục có chi phí thấp nhất
func NewCache() (*ristretto.Cache[string, domain.EventType], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, domain.EventType]{
		NumCounters: 1e7,       // number of keys to track frequency of (10M)
		MaxCost:     100 << 20, // 100MB // maximum cost of cache (100MB)
		BufferItems: 64,        // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

func NewCaches() (*ristretto.Cache[string, []domain.EventType], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, []domain.EventType]{
		NumCounters: 1e7,       // number of keys to track frequency of (10M)
		MaxCost:     100 << 20, // 100MB // maximum cost of cache (100MB)
		BufferItems: 64,        // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

func NewEventTypeUseCase(database *config.Database, contextTimeout time.Duration, eventTypeRepository eventtyperepository.IEventTypeRepository, userRepository userrepository.IUserRepository) IEventTypeUseCase {
	cache, err := NewCache()
	if err != nil {
		panic(err)
	}

	caches, err := NewCaches()
	if err != nil {
		panic(err)
	}
	return &eventTypeUseCase{cache: cache, caches: caches, database: database, contextTimeout: contextTimeout, eventTypeRepository: eventTypeRepository, userRepository: userRepository}
}

func (e eventTypeUseCase) GetByName(ctx context.Context, name string) (domain.EventType, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := e.cache.Get(name)
	if found {
		return value, nil
	}

	if name == "" {
		return domain.EventType{}, errors.New(constants.MsgInvalidInput)
	}

	data, err := e.eventTypeRepository.GetByName(ctx, name)
	if err != nil {
		return domain.EventType{}, err
	}

	e.cache.Set(name, data, 1)
	// wait for value to pass through buffers
	e.cache.Wait()

	return data, nil
}

func (e eventTypeUseCase) GetByID(ctx context.Context, id string) (domain.EventType, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := e.cache.Get(id)
	if found {
		return value, nil
	}

	eventTypeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.EventType{}, err
	}

	data, err := e.eventTypeRepository.GetByID(ctx, eventTypeID)
	if err != nil {
		return domain.EventType{}, err
	}

	e.cache.Set(id, data, 1)
	// wait for value to pass through buffers
	e.cache.Wait()

	return data, nil
}

func (e eventTypeUseCase) GetAll(ctx context.Context) ([]domain.EventType, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := e.caches.Get("typs")
	if found {
		return value, nil
	}

	data, err := e.eventTypeRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	e.caches.Set("typs", data, 1)
	// wait for value to pass through buffers
	e.caches.Wait()

	return data, nil
}

func (e eventTypeUseCase) CreateOne(ctx context.Context, eventType *domain.EventTypeInput, currentUser string) error {
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

	if err = validate_data.ValidateEventTypeInput(eventType); err != nil {
		return err
	}

	count, err := e.eventTypeRepository.CountExist(ctx, eventType.Name)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(constants.MsgAPIConflict)
	}

	eventTypeInput := domain.EventType{
		ID:   primitive.NewObjectID(),
		Name: eventType.Name,
	}

	err = e.eventTypeRepository.CreateOne(ctx, eventTypeInput)
	if err != nil {
		return err
	}

	e.caches.Clear()

	return nil
}

func (e eventTypeUseCase) UpdateOne(ctx context.Context, id string, eventType *domain.EventTypeInput, currentUser string) error {
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

	if err := validate_data.ValidateEventTypeInput(eventType); err != nil {
		return err
	}

	count, err := e.eventTypeRepository.CountExist(ctx, eventType.Name)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(constants.MsgAPIConflict)
	}

	eventTypeId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	eventTypeInput := domain.EventType{
		ID:   eventTypeId,
		Name: eventType.Name,
	}

	err = e.eventTypeRepository.UpdateOne(ctx, eventTypeInput)
	if err != nil {
		return err
	}

	e.cache.Clear()
	e.caches.Clear()

	return nil
}

func (e eventTypeUseCase) DeleteOne(ctx context.Context, id string, currentUser string) error {
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

	eventTypeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = e.eventTypeRepository.DeleteOne(ctx, eventTypeID)
	if err != nil {
		return err
	}

	e.cache.Clear()
	e.caches.Clear()

	return nil
}
