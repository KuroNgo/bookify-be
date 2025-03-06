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

func TestFindAllOrganization(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	// Mock data
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

	// Define test_e2e cases
	tests := []struct {
		name        string
		expectedErr bool
		description string
	}{
		{
			name:        "success_find_all_organizations",
			expectedErr: false,
			description: "Should successfully fetch all organizations",
		},
	}

	// Execute test_e2e cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ur.GetAll(context.Background())

			if tt.expectedErr {
				assert.Error(t, err, tt.description)
			} else {
				assert.Nil(t, err, tt.description)
			}
		})
	}
}
