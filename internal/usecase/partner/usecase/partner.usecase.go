package partner_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	partnerrepository "bookify/internal/repository/partner/repository"
	userrepository "bookify/internal/repository/user/repository"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/validate_data"
	"context"
	"errors"
	"github.com/dgraph-io/ristretto/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IPartnerUseCase interface {
	GetByID(ctx context.Context, id string) (domain.Partner, error)
	GetAll(ctx context.Context) ([]domain.Partner, error)
	CreateOne(ctx context.Context, partner *domain.PartnerInput, currentUser string) error
	UpdateOne(ctx context.Context, id string, partner *domain.PartnerInput, currentUser string) error
	DeleteOne(ctx context.Context, id string, currentUser string) error
}

type partnerUseCase struct {
	database          *config.Database
	contextTimeout    time.Duration
	partnerRepository partnerrepository.IPartnerRepository
	userRepository    userrepository.IUserRepository
	cache             *ristretto.Cache[string, domain.Partner]
	caches            *ristretto.Cache[string, []domain.Partner]
}

// NewCache Kiểm tra bộ đệm khi đã đạt đến giới hạn MaxCost
// Nếu bộ nhớ vượt quá MaxCost, Ristretto sẽ tự động xóa các mục có chi phí thấp nhất
func NewCache() (*ristretto.Cache[string, domain.Partner], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, domain.Partner]{
		NumCounters: 1e7,       // number of keys to track frequency of (10M)
		MaxCost:     100 << 20, // 100MB // maximum cost of cache (100MB)
		BufferItems: 64,        // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

// NewCaches Kiểm tra bộ đệm khi đã đạt đến giới hạn MaxCost
// Nếu bộ nhớ vượt quá MaxCost, Ristretto sẽ tự động xóa các mục có chi phí thấp nhất
func NewCaches() (*ristretto.Cache[string, []domain.Partner], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, []domain.Partner]{
		NumCounters: 1e7,       // number of keys to track frequency of (10M)
		MaxCost:     100 << 20, // 100MB // maximum cost of cache (100MB)
		BufferItems: 64,        // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

func NewPartnerUseCase(database *config.Database, contextTimeout time.Duration, partnerRepository partnerrepository.IPartnerRepository, userRepository userrepository.IUserRepository) IPartnerUseCase {
	cache, err := NewCache()
	if err != nil {
		panic(err)
	}

	caches, err := NewCaches()
	if err != nil {
		panic(err)
	}
	return &partnerUseCase{cache: cache, caches: caches, database: database, contextTimeout: contextTimeout, partnerRepository: partnerRepository, userRepository: userRepository}
}

func (p partnerUseCase) GetByID(ctx context.Context, id string) (domain.Partner, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := p.cache.Get(id)
	if found {
		return value, nil
	}

	partnerID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Partner{}, err
	}

	data, err := p.partnerRepository.GetByID(ctx, partnerID)
	if err != nil {
		return domain.Partner{}, err
	}

	p.cache.Set(id, data, 1)
	// wait for value to pass through buffers
	p.cache.Wait()

	return data, nil
}

func (p partnerUseCase) GetAll(ctx context.Context) ([]domain.Partner, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := p.caches.Get("partners")
	if found {
		return value, nil
	}

	data, err := p.partnerRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	p.caches.Set("partners", data, 1)
	// wait for value to pass through buffers
	p.caches.Wait()

	return data, nil
}

func (p partnerUseCase) CreateOne(ctx context.Context, partner *domain.PartnerInput, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	// Convert currentUser from string to primitive objectID
	userID, err := primitive.ObjectIDFromHex(currentUser)
	if err != nil {
		return err
	}

	// Handle get by id to get user data
	userData, err := p.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// from user data, check role of user
	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	if err = validate_data.ValidatePartnerInput(partner); err != nil {
		return err
	}

	count, err := p.partnerRepository.CountExist(ctx, partner.Name)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(constants.MsgAPIConflict)
	}

	partnerInput := &domain.Partner{
		ID:    primitive.NewObjectID(),
		Name:  partner.Name,
		Email: partner.Email,
		Phone: partner.Phone,
	}

	err = p.partnerRepository.CreateOne(ctx, partnerInput)
	if err != nil {
		return err
	}

	p.caches.Wait()

	return nil
}

func (p partnerUseCase) UpdateOne(ctx context.Context, id string, partner *domain.PartnerInput, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	userID, err := primitive.ObjectIDFromHex(currentUser)
	if err != nil {
		return err
	}

	userData, err := p.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	if err = validate_data.ValidatePartnerInput(partner); err != nil {
		return err
	}

	count, err := p.partnerRepository.CountExist(ctx, partner.Name)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(constants.MsgAPIConflict)
	}

	partnerID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	partnerInput := &domain.Partner{
		ID:    partnerID,
		Name:  partner.Name,
		Email: partner.Email,
		Phone: partner.Phone,
	}

	err = p.partnerRepository.UpdateOne(ctx, partnerInput)
	if err != nil {
		return err
	}

	p.caches.Wait()
	p.cache.Wait()

	return nil
}

func (p partnerUseCase) DeleteOne(ctx context.Context, id string, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	userID, err := primitive.ObjectIDFromHex(currentUser)
	if err != nil {
		return err
	}

	userData, err := p.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if userData.Role != constants.RoleSuperAdmin {
		return errors.New(constants.MsgForbidden)
	}

	partnerID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = p.partnerRepository.DeleteOne(ctx, partnerID)
	if err != nil {
		return err
	}

	p.caches.Wait()
	p.cache.Wait()

	return nil
}
