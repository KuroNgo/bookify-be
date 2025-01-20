package unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor"
	employee_repository "bookify/internal/repository/employee/repository"
	organization_repository "bookify/internal/repository/organization/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestCreateOneEmployee(t *testing.T) {
	client, database := infrastructor.SetupTestDatabase(t)
	defer infrastructor.TearDownTestDatabase(client, t)

	// Function to clear the partner collection before each test case
	clearEmployeeCollection := func() {
		err := database.Collection("employee").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear employee collection: %v", err)
		}
	}

	// Function to clear the partner collection before each test case
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

	ur := organization_repository.NewOrganizationRepository(database, "organization")
	err := ur.CreateOne(context.Background(), mockOrganization)
	assert.Nil(t, err)

	// Mock data
	mockEmployee := &domain.Employee{
		ID:             primitive.NewObjectID(),
		OrganizationID: mockOrganization.ID, // Gắn với Organization đã mock
		FirstName:      "Jane",
		LastName:       "Smith",
		JobTitle:       "Software Engineer",
		Email:          "jane.smith@techcorp.com",
	}

	mockEmployeeNil := &domain.Employee{}

	// Define test cases
	tests := []struct {
		name        string
		input       *domain.Employee
		expectedErr bool
		description string
	}{
		{
			name:        "success_create_employee",
			input:       mockEmployee,
			expectedErr: false,
			description: "Should successfully create a employee",
		},
		{
			name:        "error_create_partner_with_nil",
			input:       mockEmployeeNil,
			expectedErr: true,
			description: "Should return error when creating a employee with nil fields",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEmployeeCollection() // Clear the collection before each test

			ur := employee_repository.NewEmployeeRepository(database, "employee")
			err := ur.CreateOne(context.Background(), tt.input)

			if tt.expectedErr {
				assert.Error(t, err, tt.description)
			} else {
				assert.Nil(t, err, tt.description)
			}
		})
	}
}
