package unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor"
	organization_repository "bookify/internal/repository/organization/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestCreateOneOrganization(t *testing.T) {
	client, database := infrastructor.SetupTestDatabase(t)
	defer infrastructor.TearDownTestDatabase(client, t)

	// Function to clear the partner collection before each test case
	clearOrganizationCollection := func() {
		err := database.Collection("organization").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear organization collection: %v", err)
		}
	}

	// Mock data
	mockOrganization := &domain.Organization{
		ID:            primitive.NewObjectID(),
		Name:          "Tech Corp",
		ContactPerson: "John Doe",
		Email:         "john.doe@techcorp.com",
		Phone:         "0329245971",
	}

	mockOrganizationNil := &domain.Organization{}

	// Define test cases
	tests := []struct {
		name        string
		input       *domain.Organization
		expectedErr bool
		description string
	}{
		{
			name:        "success_create_partner",
			input:       mockOrganization,
			expectedErr: false,
			description: "Should successfully create a partner",
		},
		{
			name:        "error_create_partner_with_nil",
			input:       mockOrganizationNil,
			expectedErr: true,
			description: "Should return error when creating a partner with nil fields",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearOrganizationCollection() // Clear the collection before each test

			ur := organization_repository.NewOrganizationRepository(database, "organization")
			err := ur.CreateOne(context.Background(), tt.input)

			if tt.expectedErr {
				assert.Error(t, err, tt.description)
			} else {
				assert.Nil(t, err, tt.description)
			}
		})
	}
}
