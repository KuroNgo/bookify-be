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

func TestCountExistEmployee(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	repo := employee_repository.NewEmployeeRepository(database, "employees")

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

	err := repo.CreateOne(context.Background(), mockEmployee)
	assert.Nil(t, err, "Should successfully create an employee")

	t.Run("success_count_existing_email", func(t *testing.T) {
		count, err := repo.CountExist(context.Background(), "johndoe@example.com")
		assert.Nil(t, err, "Should successfully count existing email")
		assert.Equal(t, int64(1), count, "Should return count as 1")
	})

	t.Run("success_count_non_existing_email", func(t *testing.T) {
		count, err := repo.CountExist(context.Background(), "nonexistent@example.com")
		assert.Nil(t, err, "Should successfully count non-existing email")
		assert.Equal(t, int64(0), count, "Should return count as 0")
	})
}
