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
	"go.mongodb.org/mongo-driver/bson/primitive"
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
}

func (v venueUseCase) GetByID(ctx context.Context, id string) (domain.Venue, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	venueID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Venue{}, err
	}

	data, err := v.venueRepository.GetByID(ctx, venueID)
	if err != nil {
		return domain.Venue{}, err
	}

	return data, nil
}

func (v venueUseCase) GetAll(ctx context.Context) ([]domain.Venue, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	data, err := v.venueRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (v venueUseCase) CreateOne(ctx context.Context, venue *domain.VenueInput, currentUser string) error {
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
	//
	//count, err := p.partnerRepository.CountExist(ctx, partner.Name)
	//if err != nil {
	//	return err
	//}
	//
	//if count > 0 {
	//	return errors.New(constants.MsgAPIConflict)
	//}

	venueInput := &domain.Venue{
		ID:          primitive.NewObjectID(),
		Capacity:    venue.Capacity,
		AddressLine: venue.AddressLine,
		City:        venue.City,
		//State:       venue.State,
		Country: venue.Country,
		//PostalCode:  venue.PostalCode,
		OnlineFlat: venue.OnlineFlat,
		LinkAttend: venue.LinkAttend,
		FromAttend: venue.FromAttend,
	}

	err = v.venueRepository.CreateOne(ctx, venueInput)
	if err != nil {
		return err
	}

	return nil
}

func (v venueUseCase) UpdateOne(ctx context.Context, id string, venue *domain.VenueInput, currentUser string) error {
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
	//
	//count, err := p.partnerRepository.CountExist(ctx, partner.Name)
	//if err != nil {
	//	return err
	//}
	//
	//if count > 0 {
	//	return errors.New(constants.MsgAPIConflict)
	//}

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
		OnlineFlat: venue.OnlineFlat,
		LinkAttend: venue.LinkAttend,
		FromAttend: venue.FromAttend,
	}

	err = v.venueRepository.UpdateOne(ctx, venueInput)
	if err != nil {
		return err
	}

	return nil
}

func (v venueUseCase) DeleteOne(ctx context.Context, id string, currentUser string) error {
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

	return nil
}

func NewVenueUseCase(database *config.Database, contextTimeout time.Duration, venueRepository venue_repository.IVenueRepository, userRepository userrepository.IUserRepository) IVenueUseCase {
	return &venueUseCase{database: database, contextTimeout: contextTimeout, venueRepository: venueRepository, userRepository: userRepository}
}
