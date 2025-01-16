package event_type_repository

import (
	"bookify/internal/domain"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/validate_data"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IEventTypeRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (domain.EventType, error)
	GetAll(ctx context.Context) ([]domain.EventType, error)
	CreateOne(ctx context.Context, eventType domain.EventType) error
	UpdateOne(ctx context.Context, eventType domain.EventType) error
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
	CountExist(ctx context.Context, eventTypeName string) (int64, error)
}

type eventTypeRepository struct {
	database            *mongo.Database
	collectionEventType string
}

func (e eventTypeRepository) CountExist(ctx context.Context, eventTypeName string) (int64, error) {
	eventTypeCollection := e.database.Collection(e.collectionEventType)

	filter := bson.M{"name": eventTypeName}
	count, err := eventTypeCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (e eventTypeRepository) GetByID(ctx context.Context, id primitive.ObjectID) (domain.EventType, error) {
	eventTypeCollection := e.database.Collection(e.collectionEventType)

	filter := bson.M{"_id": id}
	var eventType domain.EventType
	if err := eventTypeCollection.FindOne(ctx, filter).Decode(&eventType); err != nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return domain.EventType{}, nil
		}
		return domain.EventType{}, err
	}

	return eventType, nil
}

func (e eventTypeRepository) GetAll(ctx context.Context) ([]domain.EventType, error) {
	eventTypeCollection := e.database.Collection(e.collectionEventType)

	filter := bson.M{}
	cursor, err := eventTypeCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var eventTypes []domain.EventType
	for cursor.Next(ctx) {
		var eventType domain.EventType
		if err = cursor.Decode(&eventType); err != nil {
			return nil, err
		}

		eventTypes = append(eventTypes, eventType)
	}

	return eventTypes, nil
}

func (e eventTypeRepository) CreateOne(ctx context.Context, eventType domain.EventType) error {
	eventTypeCollection := e.database.Collection(e.collectionEventType)

	if err := validate_data.ValidateEventType(eventType); err != nil {
		return err
	}

	filter := bson.M{"name": eventType.Name}
	count, err := eventTypeCollection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(constants.MsgAPIConflict)
	}

	_, err = eventTypeCollection.InsertOne(ctx, eventType)
	if err != nil {
		return err
	}

	return nil
}

func (e eventTypeRepository) UpdateOne(ctx context.Context, eventType domain.EventType) error {
	eventTypeCollection := e.database.Collection(e.collectionEventType)

	if err := validate_data.ValidateEventType(eventType); err != nil {
		return err
	}

	filterCount := bson.M{"name": eventType.Name}
	count, err := eventTypeCollection.CountDocuments(ctx, filterCount)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(constants.MsgAPIConflict)
	}

	filter := bson.M{"_id": eventType.ID}
	update := bson.M{"$set": bson.M{
		"name": eventType.Name,
	}}
	_, err = eventTypeCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (e eventTypeRepository) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	eventTypeCollection := e.database.Collection(e.collectionEventType)

	if id == primitive.NilObjectID {
		return errors.New(constants.MsgDataInvalidFormat)
	}

	filter := bson.M{"_id": id}
	_, err := eventTypeCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil

}

func NewEventTypeRepository(database *mongo.Database, collectionEventType string) IEventTypeRepository {
	return &eventTypeRepository{database: database, collectionEventType: collectionEventType}
}
