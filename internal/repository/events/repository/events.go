package repository

import (
	"bookify/internal/domain"
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
	GetByIDByUserID(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) (domain.Event, error)
	GetByStartTime(ctx context.Context, startTime time.Time) ([]domain.Event, error)
	GetByStartTimeByUserID(ctx context.Context, startTime time.Time, userID primitive.ObjectID) ([]domain.Event, error)
	GetByStartTimePagination(ctx context.Context, startTime time.Time, page string) ([]domain.Event, int64, int, error)
	GetByStartTimePaginationByUserID(ctx context.Context, startTime time.Time, page string, userID primitive.ObjectID) ([]domain.Event, int64, int, error)
	GetAll(ctx context.Context) ([]domain.Event, error)
	GetAllByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.Event, error)
	GetAllPagination(ctx context.Context, page string) ([]domain.Event, int64, int, error)
	GetAllPaginationByUserID(ctx context.Context, page string, userID primitive.ObjectID) ([]domain.Event, int64, int, error)
	CreateOne(ctx context.Context, event *domain.Event) error
	UpdateOne(ctx context.Context, event *domain.Event) error
	UpdateOneByUserID(ctx context.Context, event *domain.Event, userID primitive.ObjectID) error
	DeleteOne(ctx context.Context, eventID primitive.ObjectID) error
	DeleteOneByUserID(ctx context.Context, eventID primitive.ObjectID, userID primitive.ObjectID) error
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
func (e eventRepository) GetByIDByUserID(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) (domain.Event, error) {
	collectionEvent := e.database.Collection(e.collectionEvent)

	filter := bson.M{"_id": id, "user_id": userID}
	var event domain.Event
	if err := collectionEvent.FindOne(ctx, filter).Decode(&event); err != nil {
		return domain.Event{}, err
	}

	return event, nil
}

func (e eventRepository) GetByStartTime(ctx context.Context, startTime time.Time) ([]domain.Event, error) {
	collectionEvent := e.database.Collection(e.collectionEvent)

	filter := bson.M{"start_time": startTime}
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
func (e eventRepository) GetByStartTimeByUserID(ctx context.Context, startTime time.Time, userID primitive.ObjectID) ([]domain.Event, error) {
	collectionEvent := e.database.Collection(e.collectionEvent)

	filter := bson.M{"start_time": startTime, "user_id": userID}
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

	filter := bson.M{"start_time": startTime}
	pageNumber, err := strconv.Atoi(page)
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
func (e eventRepository) GetByStartTimePaginationByUserID(ctx context.Context, startTime time.Time, page string, userID primitive.ObjectID) ([]domain.Event, int64, int, error) {
	collectionEvent := e.database.Collection(e.collectionEvent)

	filter := bson.M{"start_time": startTime, "user_id": userID}
	pageNumber, err := strconv.Atoi(page)
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
func (e eventRepository) GetAllByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.Event, error) {
	collectionEvent := e.database.Collection(e.collectionEvent)

	filter := bson.M{"user_id": userID}
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
func (e eventRepository) GetAllPaginationByUserID(ctx context.Context, page string, userID primitive.ObjectID) ([]domain.Event, int64, int, error) {
	collectionEvent := e.database.Collection(e.collectionEvent)

	pageNumber, err := strconv.Atoi(page)
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

	filter := bson.M{"user_id": userID}

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

func (e eventRepository) CreateOne(ctx context.Context, event *domain.Event) error {
	collectionEvent := e.database.Collection(e.collectionEvent)

	_, err := collectionEvent.InsertOne(ctx, event)
	if err != nil {
		return err
	}

	return nil
}

func (e eventRepository) UpdateOne(ctx context.Context, event *domain.Event) error {
	collectionEvent := e.database.Collection(e.collectionEvent)

	filter := bson.M{"_id": event.ID}
	update := bson.M{"$set": bson.M{
		"name":        event.Name,
		"description": event.Description,
		"start_time":  event.StartTime,
		"end_time":    event.EndTime,
		"location":    event.Location,
		"updated_at":  time.Now(),
	}}
	_, err := collectionEvent.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
func (e eventRepository) UpdateOneByUserID(ctx context.Context, event *domain.Event, userID primitive.ObjectID) error {
	collectionEvent := e.database.Collection(e.collectionEvent)

	filter := bson.M{"_id": event.ID, "user_id": userID}
	update := bson.M{"$set": bson.M{
		"name":        event.Name,
		"description": event.Description,
		"start_time":  event.StartTime,
		"end_time":    event.EndTime,
		"location":    event.Location,
		"updated_at":  time.Now(),
	}}
	_, err := collectionEvent.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (e eventRepository) DeleteOne(ctx context.Context, eventID primitive.ObjectID) error {
	collectionEvent := e.database.Collection(e.collectionEvent)

	filter := bson.M{"_id": eventID}
	_, err := collectionEvent.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
func (e eventRepository) DeleteOneByUserID(ctx context.Context, eventID primitive.ObjectID, userID primitive.ObjectID) error {
	collectionEvent := e.database.Collection(e.collectionEvent)

	filter := bson.M{"_id": eventID, "user_id": userID}
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
