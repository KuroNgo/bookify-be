package event_wishlist_repository

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

type IEventWishlistRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (domain.EventWishlist, error)
	GetAll(ctx context.Context) ([]domain.EventWishlist, error)
	CreateOne(ctx context.Context, eventWishlist domain.EventWishlist) error
	UpdateOne(ctx context.Context, eventWishlist domain.EventWishlist) error
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
}

type eventWishlistRepository struct {
	database                *mongo.Database
	collectionEventWishlist string
}

func NewEventWishlistRepository(database *mongo.Database, collectionEventWishlist string) IEventWishlistRepository {
	return &eventWishlistRepository{database: database, collectionEventWishlist: collectionEventWishlist}
}

func (e *eventWishlistRepository) GetByID(ctx context.Context, id primitive.ObjectID) (domain.EventWishlist, error) {
	eventWishlistCollection := e.database.Collection(e.collectionEventWishlist)

	filter := bson.M{"_id": id}
	var eventWishlist domain.EventWishlist
	if err := eventWishlistCollection.FindOne(ctx, filter).Decode(&eventWishlist); err != nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return domain.EventWishlist{}, nil
		}
		return domain.EventWishlist{}, err
	}

	return eventWishlist, nil
}

func (e *eventWishlistRepository) GetAll(ctx context.Context) ([]domain.EventWishlist, error) {
	eventWishlistCollection := e.database.Collection(e.collectionEventWishlist)

	filter := bson.M{}
	cursor, err := eventWishlistCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var eventWishlists []domain.EventWishlist
	for cursor.Next(ctx) {
		var eventWishlist domain.EventWishlist
		if err = cursor.Decode(&eventWishlist); err != nil {
			return nil, err
		}

		eventWishlists = append(eventWishlists, eventWishlist)
	}

	return eventWishlists, nil
}

func (e *eventWishlistRepository) CreateOne(ctx context.Context, eventWishlist domain.EventWishlist) error {
	eventWishlistCollection := e.database.Collection(e.collectionEventWishlist)

	if err := validate_data.ValidateEventWishlist(eventWishlist); err != nil {
		return err
	}

	filter := bson.M{"event_id": eventWishlist.EventID}
	count, err := eventWishlistCollection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(constants.MsgAPIConflict)
	}

	_, err = eventWishlistCollection.InsertOne(ctx, eventWishlist)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventWishlistRepository) UpdateOne(ctx context.Context, eventWishlist domain.EventWishlist) error {
	eventWishlistCollection := e.database.Collection(e.collectionEventWishlist)

	if err := validate_data.ValidateEventWishlist(eventWishlist); err != nil {
		return err
	}

	filterCount := bson.M{"event_id": eventWishlist.EventID}
	count, err := eventWishlistCollection.CountDocuments(ctx, filterCount)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(constants.MsgAPIConflict)
	}

	filter := bson.M{"_id": eventWishlist.ID}
	update := bson.M{"$set": bson.M{
		"event_id": eventWishlist.EventID,
		"notes":    eventWishlist.Notes,
	}}
	_, err = eventWishlistCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventWishlistRepository) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	eventWishlistCollection := e.database.Collection(e.collectionEventWishlist)

	filter := bson.M{"_id": id}
	_, err := eventWishlistCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
