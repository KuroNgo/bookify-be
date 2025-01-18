package usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	event_repository "bookify/internal/repository/events/repository"
	"context"
	mongo_driven "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type IEventUseCase interface {
	GetByID(ctx context.Context, id string) (domain.EventResponse, error)
	GetByStartTime(ctx context.Context, startTime string) (domain.EventResponse, error)
	GetByStartTimePagination(ctx context.Context, startTime string, page string) (domain.EventResponse, int64, int, error)
	GetAll(ctx context.Context) (domain.EventResponse, error)
	GetAllPagination(ctx context.Context, page string) (domain.EventResponse, int64, int, error)
	CreateOne(ctx context.Context, event *domain.EventInput) error
	UpdateOne(ctx context.Context, id string, event *domain.EventInput) error
	DeleteOne(ctx context.Context, eventID string) error
}

type eventUseCase struct {
	database        *config.Database
	contextTimeout  time.Duration
	eventRepository event_repository.IEventRepository
	client          *mongo_driven.Client
}

func NewEventUseCase(database *config.Database, contextTimeout time.Duration, eventRepository event_repository.IEventRepository,
	client *mongo_driven.Client) IEventUseCase {
	return &eventUseCase{database: database, contextTimeout: contextTimeout, eventRepository: eventRepository, client: client}
}

func (e eventUseCase) GetByID(ctx context.Context, id string) (domain.EventResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (e eventUseCase) GetByStartTime(ctx context.Context, startTime string) (domain.EventResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (e eventUseCase) GetByStartTimePagination(ctx context.Context, startTime string, page string) (domain.EventResponse, int64, int, error) {
	//TODO implement me
	panic("implement me")
}

func (e eventUseCase) GetAll(ctx context.Context) (domain.EventResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (e eventUseCase) GetAllPagination(ctx context.Context, page string) (domain.EventResponse, int64, int, error) {
	//TODO implement me
	panic("implement me")
}

func (e eventUseCase) CreateOne(ctx context.Context, event *domain.EventInput) error {
	//TODO implement me
	panic("implement me")
}

func (e eventUseCase) UpdateOne(ctx context.Context, id string, event *domain.EventInput) error {
	//TODO implement me
	panic("implement me")
}

func (e eventUseCase) DeleteOne(ctx context.Context, eventID string) error {
	//TODO implement me
	panic("implement me")
}
