package event_partner_repository

import (
	"bookify/internal/domain"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IEventPartnerRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (domain.EventPartner, error)
	GetAll(ctx context.Context) ([]domain.EventPartner, error)
	CreateOne(ctx context.Context, eventPartner domain.EventPartner) error
	UpdateOne(ctx context.Context, eventPartner domain.EventPartner) error
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
}

type eventPartnerRepository struct {
	database               *mongo.Database
	collectionEventPartner string
}

func NewEventTypeRepository(database *mongo.Database, collectionEventPartner string) IEventPartnerRepository {
	return &eventPartnerRepository{database: database, collectionEventPartner: collectionEventPartner}
}

func (e *eventPartnerRepository) GetByID(ctx context.Context, id primitive.ObjectID) (domain.EventPartner, error) {
	eventPartnerCollection := e.database.Collection(e.collectionEventPartner)

	filter := bson.M{"_id": id}
	var eventPartner domain.EventPartner
	if err := eventPartnerCollection.FindOne(ctx, filter).Decode(&eventPartner); err != nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return domain.EventPartner{}, nil
		}
		return domain.EventPartner{}, err
	}

	return eventPartner, nil
}

func (e *eventPartnerRepository) GetAll(ctx context.Context) ([]domain.EventPartner, error) {
	eventPartnerCollection := e.database.Collection(e.collectionEventPartner)

	filter := bson.M{}
	cursor, err := eventPartnerCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var eventPartners []domain.EventPartner
	for cursor.Next(ctx) {
		var eventPartner domain.EventPartner
		if err = cursor.Decode(&eventPartners); err != nil {
			return nil, err
		}

		eventPartners = append(eventPartners, eventPartner)
	}

	return eventPartners, nil
}

func (e *eventPartnerRepository) CreateOne(ctx context.Context, eventPartner domain.EventPartner) error {
	eventPartnerCollection := e.database.Collection(e.collectionEventPartner)

	// Need check duplicate
	_, err := eventPartnerCollection.InsertOne(ctx, eventPartner)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventPartnerRepository) UpdateOne(ctx context.Context, eventPartner domain.EventPartner) error {
	eventPartnerCollection := e.database.Collection(e.collectionEventPartner)

	filter := bson.M{"id": eventPartner.ID}
	update := bson.M{
		"$set": bson.M{
			"event_id":   eventPartner.EventID,
			"partner_id": eventPartner.PartnerID,
			"role":       eventPartner.Role,
		},
	}

	_, err := eventPartnerCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventPartnerRepository) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	eventPartnerCollection := e.database.Collection(e.collectionEventPartner)

	filter := bson.M{"id": id}
	_, err := eventPartnerCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
