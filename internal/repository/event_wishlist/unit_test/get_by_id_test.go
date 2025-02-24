package unit_test

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor/mongodb"
	partnerrepository "bookify/internal/repository/partner/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestFindByIDPartner(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	// Function to clear the partner collection before each test case
	clearPartnerCollection := func() {
		err := database.Collection("partner").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear partner collection: %v", err)
		}
	}

	clearPartnerCollection()
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
			name:        "success_find_partner_by_id",
			inputID:     mockPartner.ID,
			expectedErr: false,
			description: "Should successfully find partner by ID",
		},
		{
			name:        "error_find_partner_by_invalid_id",
			inputID:     primitive.NilObjectID,
			expectedErr: true,
			description: "Should return error when finding partner with invalid ID",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ur.GetByID(context.Background(), tt.inputID)

			if tt.expectedErr {
				assert.Error(t, err, tt.description)
			} else {
				assert.Nil(t, err, tt.description)
			}
		})
	}
}
