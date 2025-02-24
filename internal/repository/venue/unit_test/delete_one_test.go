package unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor/mongodb"
	venuerepository "bookify/internal/repository/venue/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestDeleteOneVenue(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	// Define the mock venue
	mockVenue := &domain.Venue{
		ID:          primitive.NewObjectID(),
		Capacity:    100,
		AddressLine: "123 Main Street",
		City:        "New York",
		State:       "NY",
		Country:     "USA",
		PostalCode:  "10001",
		OnlineFlat:  false,
	}

	// Define test cases
	tests := []struct {
		name        string
		venueID     primitive.ObjectID
		expectedErr bool
		description string
	}{
		{
			name:        "success_delete_venue",
			venueID:     mockVenue.ID,
			expectedErr: false,
			description: "Should successfully delete the venue",
		},
		{
			name:        "error_delete_venue_with_nil_id",
			venueID:     primitive.NilObjectID,
			expectedErr: true,
			description: "Should return error when trying to delete with a nil ID",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear the venue collection before each test
			clearVenueCollection := func() {
				err := database.Collection("venue").Drop(context.Background())
				if err != nil {
					t.Fatalf("Failed to clear venue collection: %v", err)
				}
			}

			clearVenueCollection() // Clear collection before each test
			ur := venuerepository.NewVenueRepository(database, "venue")
			err := ur.CreateOne(context.Background(), mockVenue)
			assert.Nil(t, err)

			// Test the DeleteOne functionality
			err = ur.DeleteOne(context.Background(), tt.venueID)

			if tt.expectedErr {
				assert.Error(t, err, tt.description)
			} else {
				assert.Nil(t, err, tt.description)
			}
		})
	}
}
