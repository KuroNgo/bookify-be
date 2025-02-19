package employee_repository

import (
	"bookify/internal/domain"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/validate_data"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type IEmployeeRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (domain.Employee, error)
	GetAll(ctx context.Context) ([]domain.Employee, error)
	CreateOne(ctx context.Context, employee *domain.Employee) error
	UpdateOne(ctx context.Context, employee *domain.Employee) error
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
	DeleteSoft(ctx context.Context, id primitive.ObjectID) error
	Restore(ctx context.Context, id primitive.ObjectID) error
	CountExist(ctx context.Context, email string) (int64, error)
}

type employeeRepository struct {
	database           *mongo.Database
	collectionEmployee string
}

func NewEmployeeRepository(database *mongo.Database, collectionEmployee string) IEmployeeRepository {
	return &employeeRepository{database: database, collectionEmployee: collectionEmployee}
}

func (e employeeRepository) GetByID(ctx context.Context, id primitive.ObjectID) (domain.Employee, error) {
	employeeCollection := e.database.Collection(e.collectionEmployee)

	filter := bson.M{"_id": id, "status": "enabled"}
	var employee domain.Employee
	if err := employeeCollection.FindOne(ctx, filter).Decode(&employee); err != nil {
		return domain.Employee{}, err
	}

	return employee, nil
}

func (e employeeRepository) GetAll(ctx context.Context) ([]domain.Employee, error) {
	employeeCollection := e.database.Collection(e.collectionEmployee)

	filter := bson.M{"status": "enabled"}
	cursor, err := employeeCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var employees []domain.Employee
	for cursor.Next(ctx) {
		var employee domain.Employee
		if err = cursor.Decode(&employee); err != nil {
			return nil, err
		}

		employees = append(employees, employee)
	}

	return employees, nil
}

func (e employeeRepository) CreateOne(ctx context.Context, employee *domain.Employee) error {
	employeeCollection := e.database.Collection(e.collectionEmployee)

	if err := validate_data.ValidateEmployee(employee); err != nil {
		return err
	}

	_, err := employeeCollection.InsertOne(ctx, employee)
	if err != nil {
		return err
	}

	return nil
}

func (e employeeRepository) UpdateOne(ctx context.Context, employee *domain.Employee) error {
	employeeCollection := e.database.Collection(e.collectionEmployee)

	if err := validate_data.ValidateEmployee(employee); err != nil {
		return err
	}

	filter := bson.M{"_id": employee.ID}
	update := bson.M{"$set": bson.M{
		"organization_id": employee.OrganizationID,
		"first_name":      employee.FirstName,
		"last_name":       employee.LastName,
		"email":           employee.Email,
		"job_title":       employee.JobTitle,
		"updated_at":      employee.UpdatedAt,
	}}
	_, err := employeeCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (e employeeRepository) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	employeeCollection := e.database.Collection(e.collectionEmployee)

	if id == primitive.NilObjectID {
		return errors.New(constants.MsgInvalidInput)
	}

	filter := bson.M{"_id": id}
	_, err := employeeCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (e employeeRepository) DeleteSoft(ctx context.Context, id primitive.ObjectID) error {
	employeeCollection := e.database.Collection(e.collectionEmployee)

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"status":     "disabled",
		"updated_at": time.Now(),
	}}

	_, err := employeeCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (e employeeRepository) Restore(ctx context.Context, id primitive.ObjectID) error {
	employeeCollection := e.database.Collection(e.collectionEmployee)

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"status":     "enabled",
		"updated_at": time.Now(),
	}}
	_, err := employeeCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (e employeeRepository) CountExist(ctx context.Context, email string) (int64, error) {
	employeeCollection := e.database.Collection(e.collectionEmployee)

	filter := bson.M{"email": email}
	count, err := employeeCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}
