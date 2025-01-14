package usecase

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	"bookify/internal/repository/events/repository"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo_driven "go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

type IEventUseCase interface {
	GetByID(ctx context.Context, id string) (domain.EventResponse, error)
	GetByIDByUserID(ctx context.Context, id string, userID string) (domain.EventResponse, error)
	GetByStartTime(ctx context.Context, startTime string) (domain.EventResponse, error)
	GetByStartTimeByUserID(ctx context.Context, startTime string, userID string) (domain.EventResponse, error)
	GetByStartTimePagination(ctx context.Context, startTime string, page string) (domain.EventResponse, int64, int, error)
	GetByStartTimePaginationByUserID(ctx context.Context, startTime string, page string, userID string) (domain.EventResponse, int64, int, error)
	GetAll(ctx context.Context) (domain.EventResponse, error)
	GetAllByUserID(ctx context.Context, userID string) (domain.EventResponse, error)
	GetAllPagination(ctx context.Context, page string) (domain.EventResponse, int64, int, error)
	GetAllPaginationByUserID(ctx context.Context, page string, userID string) (domain.EventResponse, int64, int, error)
	CreateOne(ctx context.Context, event *domain.EventInput) error
	UpdateOne(ctx context.Context, id string, event *domain.EventInput) error
	DeleteOne(ctx context.Context, eventID string) error
}

type eventUseCase struct {
	database        *config.Database
	contextTimeout  time.Duration
	eventRepository repository.IEventRepository
	client          *mongo_driven.Client
}

func NewEventUseCase(database *config.Database, contextTimeout time.Duration, eventRepository repository.IEventRepository,
	client *mongo_driven.Client) IEventUseCase {
	return &eventUseCase{database: database, contextTimeout: contextTimeout, eventRepository: eventRepository, client: client}
}

func (e eventUseCase) GetByID(ctx context.Context, id string) (domain.EventResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	eventID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.EventResponse{}, fmt.Errorf("invalid ID format: %w", err)
	}

	eventData, err := e.eventRepository.GetByID(ctx, eventID)
	if err != nil {
		return domain.EventResponse{}, fmt.Errorf("failed to fetch event: %w", err)
	}

	output := domain.EventResponse{
		Event: []domain.Event{eventData},
	}

	return output, nil
}

func (e eventUseCase) GetByIDByUserID(ctx context.Context, id string, userID string) (domain.EventResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (e eventUseCase) GetByStartTime(ctx context.Context, startTime string) (domain.EventResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	layout := "2006-01-02 15:04:05"

	parsedTime, err := time.Parse(layout, startTime)
	if err != nil {
		fmt.Println("Error:", err)
	}

	eventData, err := e.eventRepository.GetByStartTime(ctx, parsedTime)
	if err != nil {
		return domain.EventResponse{}, err
	}

	output := domain.EventResponse{
		Event: eventData,
	}

	return output, nil
}

func (e eventUseCase) GetByStartTimeByUserID(ctx context.Context, startTime string, userID string) (domain.EventResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (e eventUseCase) GetByStartTimePagination(ctx context.Context, startTime string, page string) (domain.EventResponse, int64, int, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	layout := "2006-01-02 15:04:05"

	parsedTime, err := time.Parse(layout, startTime)
	if err != nil {
		fmt.Println("Error:", err)
	}

	eventData, pageOutput, currentPage, err := e.eventRepository.GetByStartTimePagination(ctx, parsedTime, page)
	if err != nil {
		return domain.EventResponse{}, 0, 0, err
	}

	output := domain.EventResponse{
		Event: eventData,
	}

	return output, pageOutput, currentPage, nil
}

func (e eventUseCase) GetByStartTimePaginationByUserID(ctx context.Context, startTime string, page string, userID string) (domain.EventResponse, int64, int, error) {
	//TODO implement me
	panic("implement me")
}

func (e eventUseCase) GetAll(ctx context.Context) (domain.EventResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	eventData, err := e.eventRepository.GetAll(ctx)
	if err != nil {
		return domain.EventResponse{}, err
	}

	output := domain.EventResponse{
		Event: eventData,
	}

	return output, nil
}

func (e eventUseCase) GetAllByUserID(ctx context.Context, userID string) (domain.EventResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (e eventUseCase) GetAllPagination(ctx context.Context, page string) (domain.EventResponse, int64, int, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	eventData, pageOutput, currentPage, err := e.eventRepository.GetAllPagination(ctx, page)
	if err != nil {
		return domain.EventResponse{}, 0, 0, err
	}

	output := domain.EventResponse{
		Event: eventData,
	}

	return output, pageOutput, currentPage, nil
}

func (e eventUseCase) GetAllPaginationByUserID(ctx context.Context, page string, userID string) (domain.EventResponse, int64, int, error) {
	//TODO implement me
	panic("implement me")
}

func (e eventUseCase) CreateOne(ctx context.Context, event *domain.EventInput) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// Validate event input
	if err := validateEvent(event); err != nil {
		return fmt.Errorf("event validation failed: %w", err)
	}

	// Parse UserID
	userID, err := primitive.ObjectIDFromHex(event.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	// Parse StartTime and EndTime
	parsedStartTime, parsedEndTime, err := parseEventTimes(event.StartTime, event.EndTime)
	if err != nil {
		return fmt.Errorf("time parsing failed: %w", err)
	}

	// Check if event already exists
	isDuplicate, err := e.checkEventDuplicate(ctx, event.Name, userID, parsedStartTime, parsedEndTime)
	if err != nil {
		return fmt.Errorf("failed to check event existence: %w", err)
	}

	if isDuplicate {
		return errors.New("the event already exists in the database")
	}

	// Map EventInput to Event
	newEvent := domain.Event{
		ID:          primitive.NewObjectID(),
		Title:       event.Name,
		Description: event.Description,
		StartTime:   parsedStartTime,
		EndTime:     parsedEndTime,
		Location:    event.Location,
		//AttendanceID: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save event to repository
	if err := e.eventRepository.CreateOne(ctx, &newEvent); err != nil {
		return fmt.Errorf("failed to save event: %w", err)
	}

	return nil
}

func parseEventTimes(startTimeStr, endTimeStr string) (time.Time, time.Time, error) {
	const layout = "2006-01-02 15:04:05"

	startTime, err := time.Parse(layout, startTimeStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start time format, must be 'YYYY-MM-DD HH:MM:SS'")
	}

	endTime, err := time.Parse(layout, endTimeStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end time format, must be 'YYYY-MM-DD HH:MM:SS'")
	}

	return startTime, endTime, nil
}

func (e eventUseCase) checkEventDuplicate(ctx context.Context, name string, userID primitive.ObjectID, startTime, endTime time.Time) (bool, error) {
	count, err := e.eventRepository.CountEventExist(ctx, name, userID, startTime, endTime)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (e eventUseCase) UpdateOne(ctx context.Context, id string, event *domain.EventInput) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// Validate event input
	if err := validateEvent(event); err != nil {
		return fmt.Errorf("event validation failed: %w", err)
	}

	// Parse Event ID
	eventID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid event ID format: %w", err)
	}

	// Parse User ID
	userID, err := primitive.ObjectIDFromHex(event.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	// Parse StartTime and EndTime
	parsedStartTime, parsedEndTime, err := parseEventTimes(event.StartTime, event.EndTime)
	if err != nil {
		return fmt.Errorf("time parsing failed: %w", err)
	}

	// Check if event already exists
	isDuplicate, err := e.checkEventDuplicate(ctx, event.Name, userID, parsedStartTime, parsedEndTime)
	if err != nil {
		return fmt.Errorf("failed to check event existence: %w", err)
	}
	if isDuplicate {
		return errors.New("the event already exists in the database")
	}

	// Map EventInput to Event
	updatedEvent := &domain.Event{
		ID:          eventID,
		Title:       event.Name,
		Description: event.Description,
		StartTime:   parsedStartTime,
		EndTime:     parsedEndTime,
		Location:    event.Location,
		//AttendanceID: userID,
		UpdatedAt: time.Now(),
	}

	// Update event in repository
	if err = e.eventRepository.UpdateOne(ctx, updatedEvent); err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return nil
}

func (e eventUseCase) DeleteOne(ctx context.Context, eventID string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	// Parse Event ID
	idEvent, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		return fmt.Errorf("invalid event ID format: %w", err)
	}

	// Check if event exists
	exists, err := e.eventRepository.CheckEventExist(ctx, idEvent)
	if err != nil {
		return fmt.Errorf("failed to check event existence: %w", err)
	}
	if !exists {
		return errors.New("event does not exist")
	}

	// Delete event
	if err := e.eventRepository.DeleteOne(ctx, idEvent); err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}

func validateEvent(input *domain.EventInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return errors.New("event name cannot be empty")
	}

	if len(input.Name) > 100 {
		return errors.New("event name cannot exceed 100 characters")
	}

	if strings.TrimSpace(input.Description) == "" {
		return errors.New("event description cannot be empty")
	}

	if len(input.Description) > 1000 {
		return errors.New("event description cannot exceed 1000 characters")
	}

	if strings.TrimSpace(input.Location) == "" {
		return errors.New("event location cannot be empty")
	}

	if len(input.Location) > 500 {
		return errors.New("event location cannot exceed 500 characters")
	}

	layout := time.RFC3339
	startTime, err := time.Parse(layout, input.StartTime)
	if err != nil {
		return errors.New("event start time must be in ISO 8601 format (e.g., 2024-05-15T09:00:00Z)")
	}

	endTime, err := time.Parse(layout, input.EndTime)
	if err != nil {
		return errors.New("event end time must be in ISO 8601 format (e.g., 2024-05-15T17:00:00Z)")
	}

	if endTime.Before(startTime) {
		return errors.New("event end time cannot be earlier than the start time")
	}

	if startTime.Before(time.Now()) {
		return errors.New("event start time cannot be in the past")
	}

	if endTime.Sub(startTime) > 24*time.Hour {
		return errors.New("event duration cannot exceed 24 hours")
	}

	return nil
}
