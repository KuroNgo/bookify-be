package event_ticket_repository

import (
	"bookify/internal/domain"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IEventTicketAssignmentRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (domain.EventTicketAssignment, error)
	GetAll(ctx context.Context) ([]domain.EventTicketAssignment, error)
	CreateOne(ctx context.Context, eventTicketAssignment domain.EventTicketAssignment) error
	UpdateOne(ctx context.Context, eventTicketAssignment domain.EventTicketAssignment) error
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
}

type eventTicketAssignmentRepository struct {
	database                        *mongo.Database
	collectionEventTicketAssignment string
}

func NewEventTicketAssignmentRepository(database *mongo.Database, collectionEventTicketAssignment string) IEventTicketAssignmentRepository {
	return &eventTicketAssignmentRepository{database: database, collectionEventTicketAssignment: collectionEventTicketAssignment}
}

func (e *eventTicketAssignmentRepository) GetByID(ctx context.Context, id primitive.ObjectID) (domain.EventTicketAssignment, error) {
	eventTicketAssignmentCollection := e.database.Collection(e.collectionEventTicketAssignment)

	filter := bson.M{"_id": id}
	var eventTicketAssignment domain.EventTicketAssignment
	if err := eventTicketAssignmentCollection.FindOne(ctx, filter).Decode(&eventTicketAssignment); err != nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return domain.EventTicketAssignment{}, nil
		}
		return domain.EventTicketAssignment{}, err
	}

	return eventTicketAssignment, nil
}

func (e *eventTicketAssignmentRepository) GetAll(ctx context.Context) ([]domain.EventTicketAssignment, error) {
	eventTicketAssignmentCollection := e.database.Collection(e.collectionEventTicketAssignment)

	filter := bson.M{}
	cursor, err := eventTicketAssignmentCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var eventTicketAssignments []domain.EventTicketAssignment
	for cursor.Next(ctx) {
		var eventTicketAssignment domain.EventTicketAssignment
		if err = cursor.Decode(&eventTicketAssignment); err != nil {
			return nil, err
		}

		eventTicketAssignments = append(eventTicketAssignments, eventTicketAssignment)
	}

	return eventTicketAssignments, nil
}

func (e *eventTicketAssignmentRepository) CreateOne(ctx context.Context, eventTicketAssignment domain.EventTicketAssignment) error {
	eventTicketAssignmentCollection := e.database.Collection(e.collectionEventTicketAssignment)

	// Need check duplicate
	_, err := eventTicketAssignmentCollection.InsertOne(ctx, eventTicketAssignment)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventTicketAssignmentRepository) UpdateOne(ctx context.Context, eventTicketAssignment domain.EventTicketAssignment) error {
	eventTicketAssignmentCollection := e.database.Collection(e.collectionEventTicketAssignment)

	filter := bson.M{"id": eventTicketAssignment.ID}
	update := bson.M{
		"$set": bson.M{
			"event_id":      eventTicketAssignment.EventID,
			"attendance_id": eventTicketAssignment.AttendanceID,
			"purchase_date": eventTicketAssignment.PurchaseDate,
			"expiry_date":   eventTicketAssignment.ExpiryDate,
			"price":         eventTicketAssignment.Price,
			"ticket_type":   eventTicketAssignment.TicketType,
			"status":        eventTicketAssignment.Status,
		},
	}

	_, err := eventTicketAssignmentCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventTicketAssignmentRepository) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	eventTicketAssignmentCollection := e.database.Collection(e.collectionEventTicketAssignment)

	filter := bson.M{"id": id}
	_, err := eventTicketAssignmentCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
