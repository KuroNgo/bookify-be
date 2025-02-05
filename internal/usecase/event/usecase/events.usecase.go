package usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	event_repository "bookify/internal/repository/event/repository"
	event_type_repository "bookify/internal/repository/event_type/repository"
	organization_repository "bookify/internal/repository/organization/repository"
	user_repository "bookify/internal/repository/user/repository"
	venue_repository "bookify/internal/repository/venue/repository"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/validate_data"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo_driven "go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"time"
)

type IEventUseCase interface {
	GetByID(ctx context.Context, id string) (domain.Event, error)
	GetByStartTime(ctx context.Context, startTime string) ([]domain.Event, error)
	GetByStartTimePagination(ctx context.Context, startTime string, page string) ([]domain.Event, int64, int, error)
	GetAll(ctx context.Context) ([]domain.Event, error)
	GetAllPagination(ctx context.Context, page string) ([]domain.Event, int64, int, error)
	CreateOne(ctx context.Context, event *domain.EventInput) error
	CreateOneAsync(ctx context.Context, event *domain.EventInput) error
	UpdateOne(ctx context.Context, id string, event *domain.EventInput) error
	DeleteOne(ctx context.Context, eventID string) error
}

type eventUseCase struct {
	database               *config.Database
	contextTimeout         time.Duration
	eventRepository        event_repository.IEventRepository
	organizationRepository organization_repository.IOrganizationRepository
	eventTypeRepository    event_type_repository.IEventTypeRepository
	venueRepository        venue_repository.IVenueRepository
	userRepository         user_repository.IUserRepository
	client                 *mongo_driven.Client
}

func NewEventUseCase(database *config.Database, contextTimeout time.Duration, eventRepository event_repository.IEventRepository,
	organizationRepository organization_repository.IOrganizationRepository, eventTypeRepository event_type_repository.IEventTypeRepository,
	venueRepository venue_repository.IVenueRepository, userRepository user_repository.IUserRepository, client *mongo_driven.Client) IEventUseCase {
	return &eventUseCase{database: database, contextTimeout: contextTimeout, eventRepository: eventRepository,
		organizationRepository: organizationRepository, eventTypeRepository: eventTypeRepository, venueRepository: venueRepository,
		userRepository: userRepository, client: client}
}

func (e eventUseCase) GetByID(ctx context.Context, id string) (domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	eventID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Event{}, err
	}

	data, err := e.eventRepository.GetByID(ctx, eventID)
	if err != nil {
		return domain.Event{}, err
	}

	return data, nil
}

func (e eventUseCase) GetByStartTime(ctx context.Context, startTime string) ([]domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// Parse thời gian từ chuỗi (ISO 8601)
	parseStartTime, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return nil, errors.New(constants.MsgInvalidInput)
	}

	// Chuyển về đầu ngày (00:00:00 UTC)
	startOfDay := time.Date(parseStartTime.Year(), parseStartTime.Month(), parseStartTime.Day(), 0, 0, 0, 0, time.UTC)

	data, err := e.eventRepository.GetByStartTime(ctx, startOfDay)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e eventUseCase) GetByStartTimePagination(ctx context.Context, startTime string, page string) ([]domain.Event, int64, int, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	layout := "20/1/2025"
	parseStartTime, err := time.Parse(layout, startTime)
	if err != nil {
		return nil, 0, 0, errors.New(constants.MsgInvalidInput)
	}

	pageChoose, err := strconv.Atoi(page)
	if err != nil {
		return nil, 0, 0, err
	}

	if pageChoose < 1 {
		return nil, 0, 0, errors.New(constants.MsgInvalidInput)
	}

	data, totalPage, pageCurrent, err := e.eventRepository.GetByStartTimePagination(ctx, parseStartTime, page)
	if err != nil {
		return nil, 0, 0, err
	}

	return data, totalPage, pageCurrent, nil
}

func (e eventUseCase) GetAll(ctx context.Context) ([]domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.eventRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e eventUseCase) GetAllPagination(ctx context.Context, page string) ([]domain.Event, int64, int, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	pageChoose, err := strconv.Atoi(page)
	if err != nil {
		return nil, 0, 0, err
	}

	if pageChoose < 1 {
		return nil, 0, 0, errors.New(constants.MsgInvalidInput)
	}

	data, totalPage, pageCurrent, err := e.eventRepository.GetAllPagination(ctx, page)
	if err != nil {
		return nil, 0, 0, err
	}

	return data, totalPage, pageCurrent, nil
}

func (e eventUseCase) CreateOne(ctx context.Context, event *domain.EventInput) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// Check validate data
	if err := validate_data.ValidateEventInput(event); err != nil {
		return err
	}

	// Check exist data
	eventTypeData, err := e.eventTypeRepository.GetByName(ctx, event.EventTypeName)
	if err != nil {
		return err
	}

	organizationID, err := primitive.ObjectIDFromHex(event.OrganizationID)
	if err != nil {
		return err
	}

	organizationData, err := e.organizationRepository.GetByID(ctx, organizationID)
	if err != nil {
		return err
	}
	if organizationData.ID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}

	//venueData, err := e.venueRepository.GetByID(ctx, event.VenueID)
	//if err != nil {
	//	return err
	//}
	//if venueData.ID == primitive.NilObjectID {
	//	return errors.New(constants.MsgInvalidInput)
	//}

	parseStartTime, err := time.Parse(time.RFC3339, event.StartTime)
	if err != nil {
		return errors.New(constants.MsgInvalidInput)
	}

	parseEndTime, err := time.Parse(time.RFC3339, event.EndTime)
	if err != nil {
		return errors.New(constants.MsgInvalidInput)
	}

	eventInput := &domain.Event{
		ID:                primitive.NewObjectID(),
		EventTypeID:       eventTypeData.ID,
		OrganizationID:    organizationID,
		Title:             event.Title,
		ShortDescription:  event.ShortDescription,
		Description:       event.Description,
		ImageURL:          event.ImageURL,
		AssetURL:          event.AssetURL,
		StartTime:         parseStartTime,
		EndTime:           parseEndTime,
		Mode:              event.Mode,
		EstimatedAttendee: event.EstimatedAttendee,
		ActualAttendee:    event.ActualAttendee,
		TotalExpenditure:  event.TotalExpenditure,
	}

	err = e.eventRepository.CreateOne(ctx, eventInput)
	if err != nil {
		return err
	}

	return nil
}

func (e eventUseCase) CreateOneAsync(ctx context.Context, event *domain.EventInput) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// Validate event and venue inputs
	if err := validate_data.ValidateEventInput(event); err != nil {
		return err
	}

	// Start MongoDB session for transaction
	session, err := e.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessionCtx mongo_driven.SessionContext) (interface{}, error) {
		// Create venue
		venueInput := &domain.Venue{
			ID:          primitive.NewObjectID(),
			Capacity:    event.Capacity,
			AddressLine: event.AddressLine,
			City:        event.City,
			Country:     event.Country,
			EventMode:   event.EventMode,
			LinkAttend:  event.LinkAttend,
			FromAttend:  event.FromAttend,
		}
		if err = e.venueRepository.CreateOne(sessionCtx, venueInput); err != nil {
			return nil, err
		}

		organizationID, err := primitive.ObjectIDFromHex(event.OrganizationID)
		if err != nil {
			return nil, err
		}

		parseStartTime, err := time.Parse(time.RFC3339, event.StartTime)
		if err != nil {
			return nil, errors.New(constants.MsgInvalidInput)
		}

		parseEndTime, err := time.Parse(time.RFC3339, event.EndTime)
		if err != nil {
			return nil, errors.New(constants.MsgInvalidInput)
		}

		eventTypeData, err := e.eventTypeRepository.GetByName(ctx, event.EventTypeName)
		if err != nil {
			return nil, errors.New(constants.MsgInvalidInput)
		}

		// Create event
		eventInput := &domain.Event{
			ID:                primitive.NewObjectID(),
			EventTypeID:       eventTypeData.ID,
			VenueID:           venueInput.ID,
			OrganizationID:    organizationID,
			Title:             event.Title,
			ShortDescription:  event.ShortDescription,
			Description:       event.Description,
			ImageURL:          event.ImageURL,
			AssetURL:          event.AssetURL,
			StartTime:         parseStartTime,
			EndTime:           parseEndTime,
			Mode:              event.Mode,
			EstimatedAttendee: event.EstimatedAttendee,
			ActualAttendee:    event.ActualAttendee,
			TotalExpenditure:  event.TotalExpenditure,
			Tags:              event.Tags,
		}
		if err := e.eventRepository.CreateOne(sessionCtx, eventInput); err != nil {
			return nil, err
		}

		return nil, nil // Successfully created event and venue
	}

	// Run transaction
	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return err
	}

	return nil // Transaction successful, no need to commit explicitly
}

func (e eventUseCase) UpdateOne(ctx context.Context, id string, event *domain.EventInput) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	eventID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	eventTypeData, err := e.eventTypeRepository.GetByName(ctx, event.EventTypeName)
	if err != nil {
		return err
	}

	organizationID, err := primitive.ObjectIDFromHex(event.OrganizationID)
	if err != nil {
		return err
	}

	parseStartTime, err := time.Parse(time.RFC3339, event.StartTime)
	if err != nil {
		return errors.New(constants.MsgInvalidInput)
	}

	parseEndTime, err := time.Parse(time.RFC3339, event.EndTime)
	if err != nil {
		return errors.New(constants.MsgInvalidInput)
	}

	eventInput := &domain.Event{
		ID:                eventID,
		EventTypeID:       eventTypeData.ID,
		OrganizationID:    organizationID,
		Title:             event.Title,
		Description:       event.Description,
		ImageURL:          event.ImageURL,
		AssetURL:          event.AssetURL,
		StartTime:         parseStartTime,
		EndTime:           parseEndTime,
		Mode:              event.Mode,
		EstimatedAttendee: event.EstimatedAttendee,
		ActualAttendee:    event.ActualAttendee,
		TotalExpenditure:  event.TotalExpenditure,
		Tags:              event.Tags,
	}

	err = e.eventRepository.UpdateOne(ctx, eventInput)
	if err != nil {
		return err
	}

	return nil
}

func (e eventUseCase) DeleteOne(ctx context.Context, eventID string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		return err
	}

	err = e.eventRepository.DeleteOne(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
