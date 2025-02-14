package organization_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	organizationrepository "bookify/internal/repository/organization/repository"
	userrepository "bookify/internal/repository/user/repository"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/validate_data"
	"context"
	"errors"
	"github.com/dgraph-io/ristretto/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IOrganizationUseCase interface {
	GetByID(ctx context.Context, id string) (domain.Organization, error)
	GetByUserID(ctx context.Context, userId string) (domain.Organization, error)
	GetAll(ctx context.Context) ([]domain.Organization, error)
	CreateOne(ctx context.Context, organization *domain.OrganizationInput, currentUser string) error
	UpdateOne(ctx context.Context, id string, organization *domain.OrganizationInput, currentUser string) error
	DeleteOne(ctx context.Context, id string, currentUser string) error
}

type organizationUseCase struct {
	database               *config.Database
	contextTimeout         time.Duration
	organizationRepository organizationrepository.IOrganizationRepository
	userRepository         userrepository.IUserRepository
	cache                  *ristretto.Cache[string, domain.Organization]
	caches                 *ristretto.Cache[string, []domain.Organization]
}

// NewCache Kiểm tra bộ đệm khi đã đạt đến giới hạn MaxCost
// Nếu bộ nhớ vượt quá MaxCost, Ristretto sẽ tự động xóa các mục có chi phí thấp nhất
func NewCache() (*ristretto.Cache[string, domain.Organization], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, domain.Organization]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M)
		MaxCost:     1 << 30, // maximum cost of cache (1GB)
		BufferItems: 64,      // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

// NewCaches Kiểm tra bộ đệm khi đã đạt đến giới hạn MaxCost
// Nếu bộ nhớ vượt quá MaxCost, Ristretto sẽ tự động xóa các mục có chi phí thấp nhất
func NewCaches() (*ristretto.Cache[string, []domain.Organization], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, []domain.Organization]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M)
		MaxCost:     1 << 30, // maximum cost of cache (1GB)
		BufferItems: 64,      // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

func NewOrganizationUseCase(database *config.Database, contextTimeout time.Duration, organizationRepository organizationrepository.IOrganizationRepository, userRepository userrepository.IUserRepository) IOrganizationUseCase {
	cache, err := NewCache()
	if err != nil {
		panic(err)
	}

	caches, err := NewCaches()
	if err != nil {
		panic(err)
	}
	return &organizationUseCase{cache: cache, caches: caches, database: database, contextTimeout: contextTimeout, organizationRepository: organizationRepository, userRepository: userRepository}
}

func (o organizationUseCase) GetByUserID(ctx context.Context, userId string) (domain.Organization, error) {
	ctx, cancel := context.WithTimeout(ctx, o.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := o.cache.Get(userId)
	if found {
		return value, nil
	}

	userID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return domain.Organization{}, err
	}

	data, err := o.organizationRepository.GetByUserID(ctx, userID)
	if err != nil {
		return domain.Organization{}, err
	}

	o.cache.Set(userId, data, 1)
	// wait for value to pass through buffers
	o.cache.Wait()

	return data, nil
}

func (o organizationUseCase) GetByID(ctx context.Context, id string) (domain.Organization, error) {
	ctx, cancel := context.WithTimeout(ctx, o.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := o.cache.Get(id)
	if found {
		return value, nil
	}

	organizationId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Organization{}, err
	}

	data, err := o.organizationRepository.GetByID(ctx, organizationId)
	if err != nil {
		return domain.Organization{}, err
	}

	o.cache.Set(id, data, 1)
	// wait for value to pass through buffers
	o.cache.Wait()

	return data, nil
}

func (o organizationUseCase) GetAll(ctx context.Context) ([]domain.Organization, error) {
	ctx, cancel := context.WithTimeout(ctx, o.contextTimeout)
	defer cancel()

	// get value from cache
	value, found := o.caches.Get("organizations")
	if found {
		return value, nil
	}

	data, err := o.organizationRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	o.caches.Set("organizations", data, 1)
	// wait for value to pass through buffers
	o.cache.Wait()

	return data, nil
}

func (o organizationUseCase) CreateOne(ctx context.Context, organization *domain.OrganizationInput, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, o.contextTimeout)
	defer cancel()

	//userID, err := primitive.ObjectIDFromHex(currentUser)
	//if err != nil {
	//	return err
	//}

	//userData, err := o.userRepository.GetByID(ctx, userID)
	//if err != nil {
	//	return err
	//}

	// Đối với organization, việc tạo organization sẽ do user đã đăng ký gói pre plan (tức hệ thống sẽ dựa trên
	// thuộc tính isPaid và planCategory) để thực hiện xây dựng)
	//if userData.Role != constants.RoleUser {
	//	return errors.New(constants.MsgForbidden)
	//}

	if err := validate_data.ValidateOrganizationInput(organization); err != nil {
		return err
	}

	count, err := o.organizationRepository.CountExist(ctx, organization.Name)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(constants.MsgAPIConflict)
	}

	organizationInput := &domain.Organization{
		ID:            primitive.NewObjectID(),
		Name:          organization.Name,
		ContactPerson: organization.ContactPerson,
		Email:         organization.Email,
		Phone:         organization.Phone,
	}

	err = o.organizationRepository.CreateOne(ctx, organizationInput)
	if err != nil {
		return err
	}

	o.caches.Clear()

	return nil
}

func (o organizationUseCase) UpdateOne(ctx context.Context, id string, organization *domain.OrganizationInput, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, o.contextTimeout)
	defer cancel()

	//userID, err := primitive.ObjectIDFromHex(currentUser)
	//if err != nil {
	//	return err
	//}

	//userData, err := o.userRepository.GetByID(ctx, userID)
	//if err != nil {
	//	return err
	//}

	// Đối với organization, việc tạo organization sẽ do user đã đăng ký gói pre plan (tức hệ thống sẽ dựa trên
	// thuộc tính isPaid và planCategory) để thực hiện xây dựng)
	//if userData.Role != constants.RoleUser {
	//	return errors.New(constants.MsgForbidden)
	//}

	if err := validate_data.ValidateOrganizationInput(organization); err != nil {
		return err
	}

	count, err := o.organizationRepository.CountExist(ctx, organization.Name)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(constants.MsgAPIConflict)
	}

	organizationId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	organizationInput := &domain.Organization{
		ID:            organizationId,
		Name:          organization.Name,
		ContactPerson: organization.ContactPerson,
		Email:         organization.Email,
		Phone:         organization.Phone,
	}

	err = o.organizationRepository.UpdateOne(ctx, organizationInput)
	if err != nil {
		return err
	}

	o.caches.Clear()
	o.cache.Clear()

	return nil
}

func (o organizationUseCase) DeleteOne(ctx context.Context, id string, currentUser string) error {
	ctx, cancel := context.WithTimeout(ctx, o.contextTimeout)
	defer cancel()

	//userID, err := primitive.ObjectIDFromHex(currentUser)
	//if err != nil {
	//	return err
	//}

	//userData, err := o.userRepository.GetByID(ctx, userID)
	//if err != nil {
	//	return err
	//}

	// Đối với organization, việc tạo organization sẽ do user đã đăng ký gói pre plan (tức hệ thống sẽ dựa trên
	// thuộc tính isPaid và planCategory) để thực hiện xây dựng)
	//if userData.Role != constants.RoleUser {
	//	return errors.New(constants.MsgForbidden)
	//}

	organizationId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = o.organizationRepository.DeleteOne(ctx, organizationId)
	if err != nil {
		return err
	}

	o.caches.Clear()
	o.cache.Clear()

	return nil
}

// Đối với organization, việc tạo organization sẽ do user đã đăng ký gói pre plan (tức hệ thống sẽ dựa trên
// thuộc tính isPaid và planCategory) để thực hiện xây dựng)
// Phân cấp plan bao gồm:  Free plan, Basic plan, Pro plan, Custom plan, Enterprise plan
// Dối với Free plan: người dùng không thể tạo Organization, nhưng thực hiện được việc tạo các event nhỏ (chỉ áp dụng online, tham gia event của người dùng khác)
