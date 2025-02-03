package venue_repository

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

type IVenueRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (domain.Venue, error)
	GetAll(ctx context.Context) ([]domain.Venue, error)
	CreateOne(ctx context.Context, venue *domain.Venue) error
	UpdateOne(ctx context.Context, venue *domain.Venue) error
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
	//CountExist(ctx context.Context, name string) (int64, error)
}

type venueRepository struct {
	database        *mongo.Database
	collectionVenue string
}

func (v venueRepository) GetByID(ctx context.Context, id primitive.ObjectID) (domain.Venue, error) {
	venueCollection := v.database.Collection(v.collectionVenue)

	filter := bson.M{"_id": id}
	var venue domain.Venue
	if err := venueCollection.FindOne(ctx, filter).Decode(&venue); err != nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return domain.Venue{}, nil
		}
		return domain.Venue{}, err
	}

	return venue, nil
}

func (v venueRepository) GetAll(ctx context.Context) ([]domain.Venue, error) {
	venueCollection := v.database.Collection(v.collectionVenue)

	filter := bson.M{}
	cursor, err := venueCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var venues []domain.Venue
	for cursor.Next(ctx) {
		var venue domain.Venue
		if err = cursor.Decode(&venue); err != nil {
			return nil, err
		}

		venues = append(venues, venue)
	}

	return venues, nil
}

func (v venueRepository) CreateOne(ctx context.Context, venue *domain.Venue) error {
	venueCollection := v.database.Collection(v.collectionVenue)

	if err := validate_data.ValidateVenue(venue); err != nil {
		return err
	}

	//filter := bson.M{"name": venue.Name}
	//count, err := organizationCollection.CountDocuments(ctx, filter)
	//if err != nil {
	//	return err
	//}
	//
	//if count > 0 {
	//	return errors.New(constants.MsgAPIConflict)
	//}

	_, err := venueCollection.InsertOne(ctx, venue)
	if err != nil {
		return err
	}

	return nil
}

func (v venueRepository) UpdateOne(ctx context.Context, venue *domain.Venue) error {
	venueCollection := v.database.Collection(v.collectionVenue)

	if err := validate_data.ValidateVenue(venue); err != nil {
		return err
	}

	filter := bson.M{"_id": venue.ID}
	update := bson.M{"$set": bson.M{
		"capacity":     venue.Capacity,
		"address_line": venue.AddressLine,
		"city":         venue.City,
		//"state":        venue.State,
		"country": venue.Country,
		//"postal_code":  venue.PostalCode,
		"online_flat": venue.OnlineFlat,
	}}
	_, err := venueCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (v venueRepository) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	venueCollection := v.database.Collection(v.collectionVenue)

	if id == primitive.NilObjectID {
		return errors.New(constants.MsgDataInvalidFormat)
	}

	filter := bson.M{"_id": id}
	_, err := venueCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

//func (v venueRepository) CountExist(ctx context.Context, name string) (int64, error) {
//	panic("implement me")
//}

func NewVenueRepository(database *mongo.Database, collectionVenue string) IVenueRepository {
	return &venueRepository{database: database, collectionVenue: collectionVenue}
}
