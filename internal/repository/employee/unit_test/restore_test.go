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

func TestRestoreEmployee(t *testing.T) {
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
		Status:         "disabled",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		WhoCreated:     "admin",
	}

	err := employeeRepo.CreateOne(context.Background(), mockEmployee)
	assert.Nil(t, err, "Should successfully create a disabled employee")

	t.Run("success_restore_employee", func(t *testing.T) {
		err := employeeRepo.Restore(context.Background(), mockEmployee.ID)
		assert.Nil(t, err, "Should successfully restore the employee")
		restoredEmployee, err := employeeRepo.GetByID(context.Background(), mockEmployee.ID)
		assert.Nil(t, err, "Should successfully retrieve the restored employee")
		assert.Equal(t, "enabled", restoredEmployee.Status, "Employee status should be 'enabled' after restore")
	})
}
