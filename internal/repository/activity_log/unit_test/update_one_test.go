package unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor/mongodb"
	employee_repository "bookify/internal/repository/employee/repository"
	organizationrepository "bookify/internal/repository/organization/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestUpdateOneEmployee(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	// Function to clear the organization collection before each test case
	clearOrganizationCollection := func() {
		err := database.Collection("organization").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear organization collection: %v", err)
		}
	}

	clearOrganizationCollection()
	mockOrganization := &domain.Organization{
		ID:            primitive.NewObjectID(),
		Name:          "Tech Corp",
		ContactPerson: "John Doe",
		Email:         "john.doe@techcorp.com",
		Phone:         "0329245971",
	}

	ur := organizationrepository.NewOrganizationRepository(database, "organization")
	err := ur.CreateOne(context.Background(), mockOrganization)
	assert.Nil(t, err)

	clearEmployeeCollection := func() {
		err := database.Collection("employee").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear employee collection: %v", err)
		}
	}

	clearEmployeeCollection()
	mockEmployee := &domain.Employee{
		ID:             primitive.NewObjectID(),
		OrganizationID: mockOrganization.ID, // Gắn với Organization đã mock
		FirstName:      "Jane",
		LastName:       "Smith",
		JobTitle:       "Software Engineer",
		Email:          "jane.smith@techcorp.com",
	}

	em := employee_repository.NewEmployeeRepository(database, "employee")
	err = em.CreateOne(context.Background(), mockEmployee)
	assert.Nil(t, err)

	// Define test cases
	tests := []struct {
		name        string
		inputData   *domain.Employee
		expectedErr bool
		description string
	}{
		{
			name: "success_update_organization",
			inputData: &domain.Employee{
				ID:             mockEmployee.ID,
				OrganizationID: mockOrganization.ID, // Gắn với Organization đã mock
				FirstName:      "Andrew",
				LastName:       "Kuro",
				JobTitle:       "Software Engineer",
				Email:          "jane.smith@techcorp.com",
			},
			expectedErr: false,
			description: "Should successfully update the organization",
		},
		{
			name: "error_update_organization_invalid_id",
			inputData: &domain.Employee{
				ID:             primitive.NilObjectID,
				OrganizationID: mockOrganization.ID, // Gắn với Organization đã mock
				FirstName:      "Andrew",
				LastName:       "Kuro",
				JobTitle:       "Software Engineer",
				Email:          "jane.smith@techcorp.com",
			},
			expectedErr: true,
			description: "Should return an error when trying to update with invalid ID",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := em.UpdateOne(context.Background(), tt.inputData)

			if tt.expectedErr {
				assert.Error(t, err, tt.description)
			} else {
				assert.Nil(t, err, tt.description)
			}
		})
	}
}
