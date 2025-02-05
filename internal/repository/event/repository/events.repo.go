package event_repository

import (
	"bookify/internal/domain"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/validate_data"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

type IEventRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (domain.Event, error)
	GetByStartTime(ctx context.Context, startTime time.Time) ([]domain.Event, error)
	GetByStartTimePagination(ctx context.Context, startTime time.Time, page string) ([]domain.Event, int64, int, error)
	GetAll(ctx context.Context) ([]domain.Event, error)
	GetAllPagination(ctx context.Context, page string) ([]domain.Event, int64, int, error)
	CreateOne(ctx context.Context, event *domain.Event) error
	UpdateOne(ctx context.Context, event *domain.Event) error
	DeleteOne(ctx context.Context, eventID primitive.ObjectID) error
	CheckEventExist(ctx context.Context, id primitive.ObjectID) (bool, error)
	CountEventExist(ctx context.Context, name string, userID primitive.ObjectID, timeStart time.Time, timeEnd time.Time) (int64, error)
}

type eventRepository struct {
	database        *mongo.Database
	collectionEvent string
}

func NewEventRepository(database *mongo.Database, collectionEvent string) IEventRepository {
	return &eventRepository{database: database, collectionEvent: collectionEvent}
}

func (e eventRepository) GetByID(ctx context.Context, id primitive.ObjectID) (domain.Event, error) {
	collectionEvent := e.database.Collection(e.collectionEvent)

	filter := bson.M{"_id": id}
	var event domain.Event
	if err := collectionEvent.FindOne(ctx, filter).Decode(&event); err != nil {
		return domain.Event{}, err
	}

	return event, nil
}

func (e eventRepository) GetByStartTime(ctx context.Context, startTime time.Time) ([]domain.Event, error) {
	collectionEvent := e.database.Collection(e.collectionEvent)

	if startTime.IsZero() {
		return nil, errors.New(constants.MsgInvalidInput)
	}

	// Tạo khoảng thời gian trong ngày
	startOfDay := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond) // 23:59:59.999999999

	// Tạo bộ lọc MongoDB để tìm trong khoảng thời gian
	filter := bson.M{
		"start_time": bson.M{
			"$gte": startOfDay, // Lớn hơn hoặc bằng 00:00:00
			"$lt":  endOfDay,   // Nhỏ hơn 23:59:59
		},
	}

	cursor, err := collectionEvent.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []domain.Event
	for cursor.Next(ctx) {
		var event domain.Event
		if err = cursor.Decode(&event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (e eventRepository) GetByStartTimePagination(ctx context.Context, startTime time.Time, page string) ([]domain.Event, int64, int, error) {
	collectionEvent := e.database.Collection(e.collectionEvent)

	if startTime.IsZero() {
		return nil, 0, 0, errors.New(constants.MsgInvalidInput)
	}

	pageNumber, err := strconv.Atoi(page)
	if pageNumber < 1 {
		return nil, 0, 0, errors.New(constants.MsgInvalidInput)
	}
	if err != nil {
		return nil, 0, 0, errors.New(constants.MsgInvalidInput)
	}
	perPage := 5
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	// Đếm tổng số lượng tài liệu trong collection
	count, err := collectionEvent.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, 0, 0, err
	}

	cal1 := count / int64(perPage)
	cal2 := count % int64(perPage)
	var cal int64 = 0
	if cal2 != 0 {
		cal = cal1 + 1
	}

	filter := bson.M{"start_time": startTime}
	cursor, err := collectionEvent.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, 0, err
	}

	var events []domain.Event
	for cursor.Next(ctx) {
		var event domain.Event
		if err = cursor.Decode(&event); err != nil {
			return nil, 0, 0, err
		}

		events = append(events, event)
	}

	return events, cal, pageNumber, nil
}

func (e eventRepository) GetAll(ctx context.Context) ([]domain.Event, error) {
	collectionEvent := e.database.Collection(e.collectionEvent)

	filter := bson.M{}
	cursor, err := collectionEvent.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []domain.Event
	for cursor.Next(ctx) {
		var event domain.Event
		if err = cursor.Decode(&event); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (e eventRepository) GetAllPagination(ctx context.Context, page string) ([]domain.Event, int64, int, error) {
	collectionEvent := e.database.Collection(e.collectionEvent)

	pageNumber, err := strconv.Atoi(page)
	if pageNumber <= 1 {
		return nil, 0, 0, errors.New(constants.MsgInvalidInput)
	}
	if err != nil {
		return nil, 0, 0, errors.New("invalid page number")
	}
	perPage := 5
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	// Đếm tổng số lượng tài liệu trong collection
	count, err := collectionEvent.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, 0, 0, err
	}

	cal1 := count / int64(perPage)
	cal2 := count % int64(perPage)
	var cal int64 = 0
	if cal2 != 0 {
		cal = cal1 + 1
	}

	cursor, err := collectionEvent.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, 0, 0, err
	}

	var events []domain.Event
	for cursor.Next(ctx) {
		var event domain.Event
		if err = cursor.Decode(&event); err != nil {
			return nil, 0, 0, err
		}

		events = append(events, event)
	}

	return events, cal, pageNumber, nil
}

func (e eventRepository) CreateOne(ctx context.Context, event *domain.Event) error {
	collectionEvent := e.database.Collection(e.collectionEvent)

	if err := validate_data.ValidateEvent(event); err != nil {
		return err
	}

	_, err := collectionEvent.InsertOne(ctx, event)
	if err != nil {
		return err
	}

	return nil
}

func (e eventRepository) UpdateOne(ctx context.Context, event *domain.Event) error {
	collectionEvent := e.database.Collection(e.collectionEvent)

	if err := validate_data.ValidateEvent(event); err != nil {
		return err
	}

	filter := bson.M{"_id": event.ID}
	update := bson.M{"$set": bson.M{
		"event_type_id":      event.EventTypeID,
		"venue_id":           event.VenueID,
		"organization_id":    event.OrganizationID,
		"title":              event.Title,
		"description":        event.Description,
		"image_url":          event.ImageURL,
		"asset_url":          event.AssetURL,
		"start_time":         event.StartTime,
		"end_time":           event.EndTime,
		"mode":               event.Mode,
		"estimated_attendee": event.EstimatedAttendee,
		"actual_attendee":    event.ActualAttendee,
		"total_expenditure":  event.TotalExpenditure,
	}}
	_, err := collectionEvent.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (e eventRepository) DeleteOne(ctx context.Context, eventID primitive.ObjectID) error {
	collectionEvent := e.database.Collection(e.collectionEvent)

	if eventID == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}

	filter := bson.M{"_id": eventID}
	_, err := collectionEvent.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (e eventRepository) CountEventExist(ctx context.Context, name string, userID primitive.ObjectID, timeStart time.Time, timeEnd time.Time) (int64, error) {
	collectionEvent := e.database.Collection(e.collectionEvent)

	filter := bson.M{
		"name":    name,
		"user_id": userID,
		"time_start": bson.M{
			"$gte": timeStart,
			"$lt":  timeEnd,
		},
	}
	count, err := collectionEvent.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (e eventRepository) CheckEventExist(ctx context.Context, id primitive.ObjectID) (bool, error) {
	collectionEvent := e.database.Collection(e.collectionEvent)
	filter := bson.M{"_id": id}

	count, err := collectionEvent.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
