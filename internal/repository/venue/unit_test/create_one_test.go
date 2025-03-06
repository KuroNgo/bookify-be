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

func TestCreateOneVenue(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	// Function to clear the venue collection before each test_e2e case
	clearVenueCollection := func() {
		err := database.Collection("venue").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear venue collection: %v", err)
		}
	}

	// Define the mock venue
	mockVenue := &domain.Venue{
		ID:          primitive.NewObjectID(),
		Capacity:    100,
		AddressLine: "123 Main Street",
		City:        "New York",
		Country:     "USA",
		EventMode:   "online",
		LinkAttend:  "",
		FromAttend:  "",
	}

	// Define test_e2e cases
	tests := []struct {
		name        string
		venue       *domain.Venue
		expectedErr bool
		description string
	}{
		{
			name:        "success_create_venue",
			venue:       mockVenue,
			expectedErr: false,
			description: "Should successfully create a venue",
		},
		{
			name:        "error_create_venue_with_nil_data",
			venue:       &domain.Venue{},
			expectedErr: true,
			description: "Should return error when creating a venue with nil data",
		},
	}

	// Execute test_e2e cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearVenueCollection() // Clear the collection before each test_e2e
			ur := venuerepository.NewVenueRepository(database, "venue")
			err := ur.CreateOne(context.Background(), tt.venue)

			if tt.expectedErr {
				assert.Error(t, err, tt.description)
			} else {
				assert.Nil(t, err, tt.description)
			}
		})
	}
}
