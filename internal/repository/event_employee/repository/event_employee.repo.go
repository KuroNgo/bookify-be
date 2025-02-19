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
	GetByEmployeeID(ctx context.Context, employeeID primitive.ObjectID) (domain.EventEmployee, error)
	GetAll(ctx context.Context) ([]domain.EventEmployee, error)
	GetIncompleteTaskPercentage(ctx context.Context, employeeID primitive.ObjectID) (domain.ResultEventUnComplete, error)
	GetCompleteTaskPercentage(ctx context.Context, employeeID primitive.ObjectID) (domain.ResultEventComplete, error)
	CreateOne(ctx context.Context, eventEmployee *domain.EventEmployee) error
	UpdateOne(ctx context.Context, eventEmployee *domain.EventEmployee) error
	CreateAndUpdateOne(ctx context.Context, eventEmployee *domain.EventEmployee) error
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

func (e *eventEmployeeRepository) GetByEmployeeID(ctx context.Context, employeeID primitive.ObjectID) (domain.EventEmployee, error) {
	eventEmployeeCollection := e.database.Collection(e.collectionEventEmployee)

	filter := bson.M{"employee_id": employeeID}
	var eventEmployee domain.EventEmployee
	if err := eventEmployeeCollection.FindOne(ctx, filter).Decode(&eventEmployee); err != nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return domain.EventEmployee{}, nil
		}
		return domain.EventEmployee{}, err
	}

	return eventEmployee, nil
}

func (e *eventEmployeeRepository) GetIncompleteTaskPercentage(ctx context.Context, employeeID primitive.ObjectID) (domain.ResultEventUnComplete, error) {
	eventEmployeeCollection := e.database.Collection(e.collectionEventEmployee)

	// Pipeline Aggregation phù hợp với struct
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"employee_id": employeeID}}}, // Lọc theo employee_id
		{{"$unwind", "$task"}},                          // Tách từng task trong mảng task[]
		{{"$group", bson.M{
			"_id":         "$employee_id",
			"total_tasks": bson.M{"$sum": 1}, // Đếm tổng task
			"incomplete_tasks": bson.M{
				"$sum": bson.M{
					"$cond": bson.A{"$task.task_completed", 0, 1}, // Đếm task chưa hoàn thành
				},
			},
		}}},
	}

	// Chạy Aggregation
	cursor, err := eventEmployeeCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return domain.ResultEventUnComplete{}, err
	}
	defer cursor.Close(ctx)

	var resultUncompleted domain.ResultEventUnComplete
	if cursor.Next(ctx) {
		if err := cursor.Decode(&resultUncompleted); err != nil {
			return domain.ResultEventUnComplete{}, err
		}
	}

	// Tránh lỗi chia cho 0
	if resultUncompleted.TotalTasks == 0 {
		return domain.ResultEventUnComplete{}, nil
	}

	// Tính phần trăm số task chưa hoàn thành
	percentage := (float64(resultUncompleted.IncompleteTasks) / float64(resultUncompleted.TotalTasks)) * 100
	response := domain.ResultEventUnComplete{
		IncompleteTasks: resultUncompleted.IncompleteTasks,
		TotalTasks:      resultUncompleted.TotalTasks,
		Percentage:      percentage,
	}

	return response, nil

}

func (e *eventEmployeeRepository) GetCompleteTaskPercentage(ctx context.Context, employeeID primitive.ObjectID) (domain.ResultEventComplete, error) {
	eventEmployeeCollection := e.database.Collection(e.collectionEventEmployee)

	// Pipeline để tính phần trăm task đã hoàn thành
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"employee_id": employeeID}}}, // Lọc theo employee_id
		{{"$unwind", "$task"}},                          // Tách từng task trong danh sách task[]
		{{"$group", bson.M{
			"_id":         "$employee_id",
			"total_tasks": bson.M{"$sum": 1}, // Tổng số task
			"completed_tasks": bson.M{
				"$sum": bson.M{
					"$cond": bson.A{"$task.task_completed", 1, 0}, // Đếm task hoàn thành
				},
			},
		}}},
	}

	cursor, err := eventEmployeeCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return domain.ResultEventComplete{}, err
	}
	defer cursor.Close(ctx)

	var resultCompleted domain.ResultEventComplete
	if cursor.Next(ctx) {
		if err := cursor.Decode(&resultCompleted); err != nil {
			return domain.ResultEventComplete{}, err
		}
	}

	// Tránh lỗi chia cho 0
	if resultCompleted.TotalTasks == 0 {
		return domain.ResultEventComplete{}, nil
	}

	// Tính phần trăm số task đã hoàn thành
	percentage := (float64(resultCompleted.CompleteTasks) / float64(resultCompleted.TotalTasks)) * 100
	response := domain.ResultEventComplete{
		CompleteTasks: resultCompleted.CompleteTasks,
		TotalTasks:    resultCompleted.TotalTasks,
		Percentage:    percentage,
	}

	return response, nil
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
	update := bson.M{
		"$set": bson.M{
			"event_id":    eventEmployee.EventID,
			"employee_id": eventEmployee.EmployeeID,
		},
		"$push": bson.M{
			"task": bson.M{
				"$each": eventEmployee.Task, // Nếu có nhiều task thì dùng $each
			},
		},
	}

	_, err := eventEmployeeCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventEmployeeRepository) CreateAndUpdateOne(ctx context.Context, eventEmployee *domain.EventEmployee) error {
	//TODO implement me
	panic("implement me")
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
