package unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor/mongodb"
	partnerrepository "bookify/internal/repository/partner/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestCreateOnePartner(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	// Function to clear the partner collection before each test case
	clearPartnerCollection := func() {
		err := database.Collection("partner").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear partner collection: %v", err)
		}
	}

	// Mock data
	mockPartner := &domain.Partner{
		ID:    primitive.NewObjectID(),
		Name:  "kuro",
		Email: "kuro@gmail.com",
		Phone: "0329245971",
	}

	mockPartnerNil := &domain.Partner{}

	// Define test cases
	tests := []struct {
		name        string
		input       *domain.Partner
		expectedErr bool
		description string
	}{
		{
			name:        "success_create_partner",
			input:       mockPartner,
			expectedErr: false,
			description: "Should successfully create a partner",
		},
		{
			name:        "error_create_partner_with_nil",
			input:       mockPartnerNil,
			expectedErr: true,
			description: "Should return error when creating a partner with nil fields",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearPartnerCollection() // Clear the collection before each test

			ur := partnerrepository.NewPartnerRepository(database, "partner")
			err := ur.CreateOne(context.Background(), tt.input)

			if tt.expectedErr {
				assert.Error(t, err, tt.description)
			} else {
				assert.Nil(t, err, tt.description)
			}
		})
	}
}
