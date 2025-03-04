package venue_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	userrepository "bookify/internal/repository/user/repository"
	venue_repository "bookify/internal/repository/venue/repository"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/validate_data"
	"context"
	"errors"
	"github.com/dgraph-io/ristretto/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"time"
)

type IVenueUseCase interface {
	GetByID(ctx context.Context, id string) (domain.Venue, error)
	GetAll(ctx context.Context) ([]domain.Venue, error)
	CreateOne(ctx context.Context, venue *domain.VenueInput, currentUser string) error
	UpdateOne(ctx context.Context, id string, venue *domain.VenueInput, currentUser string) error
	DeleteOne(ctx context.Context, id string, currentUser string) error
}

type venueUseCase struct {
	database        *config.Database
	contextTimeout  time.Duration
	venueRepository venue_repository.IVenueRepository
	userRepository  userrepository.IUserRepository
	mu              *sync.Mutex
	cache           *ristretto.Cache[string, domain.Venue]
	cacheVenues     *ristretto.Cache[string, []domain.Venue]
}

// NewCache Kiểm tra bộ đệm khi đã đạt đến giới hạn MaxCost
// Nếu bộ nhớ vượt quá MaxCost, Ristretto sẽ tự động xóa các mục có chi phí thấp nhất
func NewCache() (*ristretto.Cache[string, domain.Venue], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, domain.Venue]{
		NumCounters: 1e7,       // number of keys to track frequency of (10M)
		MaxCost:     100 << 20, // 100MB // maximum cost of cache (100MB)
		BufferItems: 64,        // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

// NewCacheVenue Kiểm tra bộ đệm khi đã đạt đến giới hạn MaxCost
// Nếu bộ nhớ vượt quá MaxCost, Ristretto sẽ tự động xóa các mục có chi phí thấp nhất
func NewCacheVenue() (*ristretto.Cache[string, []domain.Venue], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, []domain.Venue]{
		NumCounters: 1e7,       // number of keys to track frequency of (10M)
		MaxCost:     100 << 20, // 100MB // maximum cost of cache (100MB)
		BufferItems: 64,        // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

func NewVenueUseCase(database *config.Database, contextTimeout time.Duration, venueRepository venue_repository.IVenueRepository, userRepository userrepository.IUserRepository) IVenueUseCase {
	cache, err := NewCache()
	if err != nil {
		panic(err)
	}

	caches, err := NewCacheVenue()
	if err != nil {
		panic(err)
	}
	return &venueUseCase{cache: cache, cacheVenues: caches, database: database, contextTimeout: contextTimeout, venueRepository: venueRepository, userRepository: userRepository}
}

func (v *venueUseCase) GetByID(ctx context.Context, id string) (domain.Venue, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := v.cache.Get(id)
	if found {
		return value, nil
	}

	venueID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Venue{}, err
	}

	data, err := v.venueRepository.GetByID(ctx, venueID)
	if err != nil {
		return domain.Venue{}, err
	}

	v.cache.Set(id, data, 1)
	// wait for value to pass through buffers
	v.cache.Wait()
	return data, nil
}

func (v *venueUseCase) GetAll(ctx context.Context) ([]domain.Venue, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := v.cacheVenues.Get("venues")
	if found {
		return value, nil
	}

	data, err := v.venueRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	v.cacheVenues.Set("venues", data, 1)
	// wait for value to pass through buffers
	v.cacheVenues.Wait()
	return data, nil
}

func (v *venueUseCase) CreateOne(ctx context.Context, venue *domain.VenueInput, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	// Convert currentUser from string to primitive objectID
	userID, err := primitive.ObjectIDFromHex(currentUser)
	if err != nil {
		return err
	}

	// Handle get by id to get user data
	userData, err := v.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// from user data, check role of user
	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	if err = validate_data.ValidateVenueInput(venue); err != nil {
		return err
	}

	venueInput := &domain.Venue{
		ID:          primitive.NewObjectID(),
		Capacity:    venue.Capacity,
		AddressLine: venue.AddressLine,
		City:        venue.City,
		//State:       venue.State,
		Country: venue.Country,
		//PostalCode:  venue.PostalCode,
		EventMode:  venue.EventMode,
		LinkAttend: venue.LinkAttend,
		FromAttend: venue.FromAttend,
	}

	err = v.venueRepository.CreateOne(ctx, venueInput)
	if err != nil {
		return err
	}

	v.cache.Clear()
	return nil
}

func (v *venueUseCase) UpdateOne(ctx context.Context, id string, venue *domain.VenueInput, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	userID, err := primitive.ObjectIDFromHex(currentUser)
	if err != nil {
		return err
	}

	userData, err := v.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	if err = validate_data.ValidateVenueInput(venue); err != nil {
		return err
	}

	venueID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	venueInput := &domain.Venue{
		ID:          venueID,
		Capacity:    venue.Capacity,
		AddressLine: venue.AddressLine,
		City:        venue.City,
		//State:       venue.State,
		Country: venue.Country,
		//PostalCode:  venue.PostalCode,
		EventMode:  venue.EventMode,
		LinkAttend: venue.LinkAttend,
		FromAttend: venue.FromAttend,
	}

	err = v.venueRepository.UpdateOne(ctx, venueInput)
	if err != nil {
		return err
	}

	v.cache.Clear()
	v.cacheVenues.Clear()

	return nil
}

func (v *venueUseCase) DeleteOne(ctx context.Context, id string, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	userID, err := primitive.ObjectIDFromHex(currentUser)
	if err != nil {
		return err
	}

	userData, err := v.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	venueID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = v.venueRepository.DeleteOne(ctx, venueID)
	if err != nil {
		return err
	}

	v.cache.Clear()
	v.cacheVenues.Clear()

	return nil
}
