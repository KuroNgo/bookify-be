package unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor"
	partner_repository "bookify/internal/repository/partner/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestUpdateOnePartner(t *testing.T) {
	client, database := infrastructor.SetupTestDatabase(t)
	defer infrastructor.TearDownTestDatabase(client, t)

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

	par := partner_repository.NewPartnerRepository(database, "partner")
	err := par.CreateOne(context.Background(), mockPartner)
	assert.Nil(t, err)

	// Define test cases
	tests := []struct {
		name         string
		inputPartner *domain.Partner
		expectedErr  bool
		description  string
	}{
		{
			name: "success_update_partner",
			inputPartner: &domain.Partner{
				ID:    mockPartner.ID,
				Name:  "andrew",
				Email: "andrew@gmail.com",
				Phone: "0329245971",
			},
			expectedErr: false,
			description: "Should successfully update the partner",
		},
		{
			name: "error_update_partner_with_invalid_id",
			inputPartner: &domain.Partner{
				ID:    primitive.NilObjectID,
				Name:  "andrew",
				Email: "andrew@gmail.com",
				Phone: "0329245971",
			},
			expectedErr: true,
			description: "Should return error when updating partner with invalid ID",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := par.UpdateOne(context.Background(), tt.inputPartner)

			if tt.expectedErr {
				assert.Error(t, err, tt.description)
			} else {
				assert.Nil(t, err, tt.description)
			}
		})
	}
}
