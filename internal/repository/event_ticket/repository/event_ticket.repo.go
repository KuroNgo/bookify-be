package event_ticket_repository

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

type IEventTicketRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (domain.EventTicket, error)
	GetByEventID(ctx context.Context, eventId primitive.ObjectID) (domain.EventTicket, error)
	GetAll(ctx context.Context) ([]domain.EventTicket, error)
	CreateOne(ctx context.Context, eventTicket domain.EventTicket) error
	UpdateOne(ctx context.Context, eventTicket domain.EventTicket) error
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
}

type eventTicketRepository struct {
	database              *mongo.Database
	collectionEventTicket string
}

func NewEventTicketRepository(database *mongo.Database, collectionEventTicket string) IEventTicketRepository {
	return &eventTicketRepository{database: database, collectionEventTicket: collectionEventTicket}
}

func (e *eventTicketRepository) GetByID(ctx context.Context, id primitive.ObjectID) (domain.EventTicket, error) {
	collectionEventTicket := e.database.Collection(e.collectionEventTicket)

	if id == primitive.NilObjectID {
		return domain.EventTicket{}, errors.New(constants.MsgInvalidInput)
	}

	filter := bson.M{"_id": id}
	var eventTicket domain.EventTicket
	if err := collectionEventTicket.FindOne(ctx, filter).Decode(&eventTicket); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.EventTicket{}, nil
		}
		return domain.EventTicket{}, err
	}

	return eventTicket, nil
}

func (e *eventTicketRepository) GetByEventID(ctx context.Context, eventId primitive.ObjectID) (domain.EventTicket, error) {
	collectionEventTicket := e.database.Collection(e.collectionEventTicket)

	if eventId == primitive.NilObjectID {
		return domain.EventTicket{}, errors.New(constants.MsgInvalidInput)
	}

	filter := bson.M{"event_id": eventId}
	var eventTicket domain.EventTicket
	if err := collectionEventTicket.FindOne(ctx, filter).Decode(&eventTicket); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.EventTicket{}, nil
		}
		return domain.EventTicket{}, err
	}

	return eventTicket, nil
}

func (e *eventTicketRepository) GetAll(ctx context.Context) ([]domain.EventTicket, error) {
	collectionEventTicket := e.database.Collection(e.collectionEventTicket)

	filter := bson.M{}
	cursor, err := collectionEventTicket.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var eventTickets []domain.EventTicket
	for cursor.Next(ctx) {
		var eventTicket domain.EventTicket
		if err = cursor.Decode(&eventTicket); err != nil {
			return nil, err
		}

		eventTickets = append(eventTickets, eventTicket)
	}

	return eventTickets, nil
}

func (e *eventTicketRepository) CreateOne(ctx context.Context, eventTicket domain.EventTicket) error {
	collectionEventTicket := e.database.Collection(e.collectionEventTicket)

	if err := validate_data.ValidateEventTicket(eventTicket); err != nil {
		return err
	}

	_, err := collectionEventTicket.InsertOne(ctx, eventTicket)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventTicketRepository) UpdateOne(ctx context.Context, eventTicket domain.EventTicket) error {
	collectionEventTicket := e.database.Collection(e.collectionEventTicket)

	if err := validate_data.ValidateEventTicket(eventTicket); err != nil {
		return err
	}

	filter := bson.M{"_id": eventTicket.ID}
	update := bson.M{"$set": bson.M{
		"event_id": eventTicket.EventID,
		"price":    eventTicket.Price,
		"quantity": eventTicket.Quantity,
		"status":   eventTicket.Status,
	}}

	_, err := collectionEventTicket.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventTicketRepository) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	collectionEventTicket := e.database.Collection(e.collectionEventTicket)

	if id == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}

	filter := bson.M{"_id": id}
	_, err := collectionEventTicket.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
