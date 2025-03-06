package unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor/mongodb"
	employee_repository "bookify/internal/repository/employee/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestDeleteSoftEmployee(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	employeeRepo := employee_repository.NewEmployeeRepository(database, "employees")

	clearEmployeeCollection := func() {
		err := database.Collection("employees").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear employee collection: %v", err)
		}
	}

	clearEmployeeCollection()

	mockEmployee := &domain.Employee{
		ID:             primitive.NewObjectID(),
		OrganizationID: primitive.NewObjectID(),
		FirstName:      "John",
		LastName:       "Doe",
		JobTitle:       "Software Engineer",
		Email:          "johndoe@example.com",
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		WhoCreated:     "admin",
	}

	err := employeeRepo.CreateOne(context.Background(), mockEmployee)
	assert.Nil(t, err, "Should successfully create an employee")

	t.Run("success_delete_employee", func(t *testing.T) {
		err := employeeRepo.DeleteSoft(context.Background(), mockEmployee.ID)
		assert.Nil(t, err, "Should successfully soft delete the employee")
		updatedEmployee, err := employeeRepo.GetByID(context.Background(), mockEmployee.ID)
		assert.Nil(t, err, "Should successfully retrieve the employee")
		assert.Equal(t, "disabled", updatedEmployee.Status, "Employee status should be 'disabled'")
	})

	t.Run("error_delete_invalid_employee", func(t *testing.T) {
		err := employeeRepo.DeleteSoft(context.Background(), primitive.NilObjectID)
		assert.Error(t, err, "Should return an error when trying to delete with an invalid ID")
	})
}
