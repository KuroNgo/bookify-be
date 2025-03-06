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

func TestCountExist(t *testing.T) {
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
	// Mock data
	mockPartner := &domain.Partner{
		ID:    primitive.NewObjectID(),
		Name:  "kuro",
		Email: "kuro@gmail.com",
		Phone: "0329245971",
	}

	par := partnerrepository.NewPartnerRepository(database, "partner")
	err := par.CreateOne(context.Background(), mockPartner)
	assert.Nil(t, err)

	// Define test cases
	tests := []struct {
		name        string
		partnerName string
		expectedCnt int64
		expectedErr bool
		description string
	}{
		{
			name:        "success_count_existing_partner",
			partnerName: "kuro",
			expectedCnt: 1,
			expectedErr: false,
			description: "Should return count 1 for existing partner",
		},
		{
			name:        "success_count_non_existing_partner",
			partnerName: "nonexistent",
			expectedCnt: 0,
			expectedErr: false,
			description: "Should return count 0 for non-existing partner",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := par.CountExist(context.Background(), tt.partnerName)

			if tt.expectedErr {
				assert.Error(t, err, tt.description)
			} else {
				assert.Nil(t, err, tt.description)
				assert.Equal(t, tt.expectedCnt, count, "Count should match expected value")
			}
		})
	}
}
