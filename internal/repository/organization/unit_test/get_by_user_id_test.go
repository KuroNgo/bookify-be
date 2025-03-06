package unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor/mongodb"
	organizationrepository "bookify/internal/repository/organization/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestGetByUserID(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	// Function to clear the partner collection before each test case
	clearOrganizationCollection := func() {
		err := database.Collection("organization").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear organization collection: %v", err)
		}
	}
	clearOrganizationCollection()

	// Mock data
	mockUserID := primitive.NewObjectID()
	mockOrganization := &domain.Organization{
		ID:            primitive.NewObjectID(),
		UserID:        mockUserID,
		Name:          "Tech Corp",
		ContactPerson: "John Doe",
		Email:         "john.doe@techcorp.com",
		Phone:         "0329245971",
	}

	orgRepo := organizationrepository.NewOrganizationRepository(database, "organization")
	err := orgRepo.CreateOne(context.Background(), mockOrganization)
	assert.Nil(t, err)

	// Define test cases
	tests := []struct {
		name        string
		userID      primitive.ObjectID
		expectedErr bool
		description string
	}{
		{
			name:        "success_get_existing_organization",
			userID:      mockUserID,
			expectedErr: false,
			description: "Should return the organization for a valid user ID",
		},
		{
			name:        "error_get_non_existing_organization",
			userID:      primitive.NewObjectID(),
			expectedErr: true,
			description: "Should return an error for a non-existing user ID",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			org, err := orgRepo.GetByUserID(context.Background(), tt.userID)

			if tt.expectedErr {
				assert.Error(t, err, tt.description)
			} else {
				assert.Nil(t, err, tt.description)
				assert.Equal(t, mockOrganization.ID, org.ID, "Organization ID should match")
			}
		})
	}
}
