package unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor"
	partnerrepository "bookify/internal/repository/partner/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestDeleteOnePartner(t *testing.T) {
	client, database := infrastructor.SetupTestDatabase(t)
	defer infrastructor.TearDownTestDatabase(client, t)

	// Mock data
	mockPartner := &domain.Partner{
		ID:    primitive.NewObjectID(),
		Name:  "kuro",
		Email: "kuro@gmail.com",
		Phone: "0329245971",
	}

	ur := partnerrepository.NewPartnerRepository(database, "partner")
	err := ur.CreateOne(context.Background(), mockPartner)
	assert.Nil(t, err)

	// Define test cases
	tests := []struct {
		name        string
		inputID     primitive.ObjectID
		expectedErr bool
		description string
	}{
		{
			name:        "success_delete_partner",
			inputID:     mockPartner.ID,
			expectedErr: false,
			description: "Should successfully delete the partner",
		},
		{
			name:        "error_delete_invalid_partner",
			inputID:     primitive.NilObjectID,
			expectedErr: true,
			description: "Should return an error when trying to delete with invalid ID",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ur.DeleteOne(context.Background(), tt.inputID)

			if tt.expectedErr {
				assert.Error(t, err, tt.description)
			} else {
				assert.Nil(t, err, tt.description)
			}
		})
	}
}
