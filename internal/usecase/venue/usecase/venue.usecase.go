package venue_usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	userrepository "bookify/internal/repository/user/repository"
	venue_repository "bookify/internal/repository/venue/repository"
	"context"
	"time"
)

type IVenueUseCase interface {
	GetByID(ctx context.Context, id string) (domain.Venue, error)
	GetAll(ctx context.Context) ([]domain.Partner, error)
	CreateOne(ctx context.Context, venue *domain.VenueInput, currentUser string) error
	UpdateOne(ctx context.Context, venue *domain.VenueInput, currentUser string) error
	DeleteOne(ctx context.Context, id string, currentUser string) error
}

type venueUseCase struct {
	database        *config.Database
	contextTimeout  time.Duration
	venueRepository venue_repository.IVenueRepository
	userRepository  userrepository.IUserRepository
}

func (v venueUseCase) GetByID(ctx context.Context, id string) (domain.Venue, error) {
	//TODO implement me
	panic("implement me")
}

func (v venueUseCase) GetAll(ctx context.Context) ([]domain.Partner, error) {
	//TODO implement me
	panic("implement me")
}

func (v venueUseCase) CreateOne(ctx context.Context, venue *domain.VenueInput, currentUser string) error {
	//TODO implement me
	panic("implement me")
}

func (v venueUseCase) UpdateOne(ctx context.Context, venue *domain.VenueInput, currentUser string) error {
	//TODO implement me
	panic("implement me")
}

func (v venueUseCase) DeleteOne(ctx context.Context, id string, currentUser string) error {
	//TODO implement me
	panic("implement me")
}

func NewPartnerUseCase(database *config.Database, contextTimeout time.Duration, venueRepository venue_repository.IVenueRepository, userRepository userrepository.IUserRepository) IVenueUseCase {
	return &venueUseCase{database: database, contextTimeout: contextTimeout, venueRepository: venueRepository, userRepository: userRepository}
}
