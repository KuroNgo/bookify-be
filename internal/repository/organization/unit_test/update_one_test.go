package unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor"
	organizationrepository "bookify/internal/repository/organization/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestUpdateOneOrganization(t *testing.T) {
	client, database := infrastructor.SetupTestDatabase(t)
	defer infrastructor.TearDownTestDatabase(client, t)

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

	// Define test cases
	tests := []struct {
		name        string
		inputData   *domain.Organization
		expectedErr bool
		description string
	}{
		{
			name: "success_update_organization",
			inputData: &domain.Organization{
				ID:            mockOrganization.ID,
				Name:          "Tech Corp Updated",
				ContactPerson: "Jane Smith",
				Email:         "jane.smith@techcorp.com",
				Phone:         "0329245971",
			},
			expectedErr: false,
			description: "Should successfully update the organization",
		},
		{
			name: "error_update_organization_invalid_id",
			inputData: &domain.Organization{
				ID:            primitive.NilObjectID,
				Name:          "Invalid Organization",
				ContactPerson: "Invalid Person",
				Email:         "invalid@invalid.com",
				Phone:         "000-000-0000",
			},
			expectedErr: true,
			description: "Should return an error when trying to update with invalid ID",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ur.UpdateOne(context.Background(), tt.inputData)

			if tt.expectedErr {
				assert.Error(t, err, tt.description)
			} else {
				assert.Nil(t, err, tt.description)
			}
		})
	}
}
