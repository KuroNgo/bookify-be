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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IPartnerUseCase interface {
	GetByID(ctx context.Context, id string) (domain.Partner, error)
	GetAll(ctx context.Context) ([]domain.Partner, error)
	CreateOne(ctx context.Context, partner *domain.PartnerInput, currentUser string) error
	UpdateOne(ctx context.Context, partner *domain.PartnerInput, currentUser string) error
	DeleteOne(ctx context.Context, id string, currentUser string) error
}

type partnerUseCase struct {
	database          *config.Database
	contextTimeout    time.Duration
	partnerRepository partnerrepository.IPartnerRepository
	userRepository    userrepository.IUserRepository
}

func (p partnerUseCase) GetByID(ctx context.Context, id string) (domain.Partner, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	partnerID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Partner{}, err
	}

	data, err := p.partnerRepository.GetByID(ctx, partnerID)
	if err != nil {
		return domain.Partner{}, err
	}

	return data, nil
}

func (p partnerUseCase) GetAll(ctx context.Context) ([]domain.Partner, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	data, err := p.partnerRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

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

	return nil
}

func (p partnerUseCase) UpdateOne(ctx context.Context, partner *domain.PartnerInput, currentUser string) error {
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

	partnerInput := &domain.Partner{
		ID:    primitive.NewObjectID(),
		Name:  partner.Name,
		Email: partner.Email,
		Phone: partner.Phone,
	}

	err = p.partnerRepository.UpdateOne(ctx, partnerInput)
	if err != nil {
		return err
	}

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

	return nil
}

func NewPartnerUseCase(database *config.Database, contextTimeout time.Duration, partnerRepository partnerrepository.IPartnerRepository, userRepository userrepository.IUserRepository) IPartnerUseCase {
	return &partnerUseCase{database: database, contextTimeout: contextTimeout, partnerRepository: partnerRepository, userRepository: userRepository}
}
