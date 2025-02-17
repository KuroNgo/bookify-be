package event_employee_repository

import (
	"bookify/internal/domain"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IEventEmployeeRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (domain.EventEmployee, error)
	GetAll(ctx context.Context) ([]domain.EventEmployee, error)
	CreateOne(ctx context.Context, eventEmployee *domain.EventEmployee) error
	UpdateOne(ctx context.Context, eventEmployee *domain.EventEmployee) error
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
}

type eventEmployeeRepository struct {
	database                *mongo.Database
	collectionEventEmployee string
}

func NewEventEmployeeRepository(database *mongo.Database, collectionEventEmployee string) IEventEmployeeRepository {
	return &eventEmployeeRepository{database: database, collectionEventEmployee: collectionEventEmployee}
}

func (e *eventEmployeeRepository) GetByID(ctx context.Context, id primitive.ObjectID) (domain.EventEmployee, error) {
	eventEmployeeCollection := e.database.Collection(e.collectionEventEmployee)

	filter := bson.M{"_id": id}
	var eventEmployee domain.EventEmployee
	if err := eventEmployeeCollection.FindOne(ctx, filter).Decode(&eventEmployee); err != nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return domain.EventEmployee{}, nil
		}
		return domain.EventEmployee{}, err
	}

	return eventEmployee, nil
}

func (e *eventEmployeeRepository) GetAll(ctx context.Context) ([]domain.EventEmployee, error) {
	eventEmployeeCollection := e.database.Collection(e.collectionEventEmployee)

	filter := bson.M{}
	cursor, err := eventEmployeeCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var eventEmployees []domain.EventEmployee
	for cursor.Next(ctx) {
		var eventEmployee domain.EventEmployee
		if err = cursor.Decode(&eventEmployee); err != nil {
			return nil, err
		}

		eventEmployees = append(eventEmployees, eventEmployee)
	}

	return eventEmployees, nil
}

func (e *eventEmployeeRepository) CreateOne(ctx context.Context, eventEmployee *domain.EventEmployee) error {
	eventEmployeeCollection := e.database.Collection(e.collectionEventEmployee)

	//if err := validate_data.ValidateEventEmployee(eventEmployee); err != nil {
	//	return err
	//}

	_, err := eventEmployeeCollection.InsertOne(ctx, eventEmployee)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventEmployeeRepository) UpdateOne(ctx context.Context, eventEmployee *domain.EventEmployee) error {
	eventEmployeeCollection := e.database.Collection(e.collectionEventEmployee)

	//if err := validate_data.ValidateEventType(eventType); err != nil {
	//	return err
	//}

	filter := bson.M{"_id": eventEmployee.ID}
	update := bson.M{"$set": bson.M{
		"event_id":       eventEmployee.EventID,
		"employee_id":    eventEmployee.EmployeeID,
		"task":           eventEmployee.Task,
		"start_date":     eventEmployee.StartDate,
		"deadline":       eventEmployee.Deadline,
		"task_completed": eventEmployee.TaskCompleted,
	}}

	_, err := eventEmployeeCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventEmployeeRepository) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	eventEmployeeCollection := e.database.Collection(e.collectionEventEmployee)

	filter := bson.M{"_id": id}
	_, err := eventEmployeeCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
