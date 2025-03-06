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

func TestFindAllVenue(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

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

	// Define test cases
	tests := []struct {
		name        string
		venue       *domain.Venue
		expectedErr bool
	}{
		{
			name:        "success_find_all_venue",
			venue:       mockVenue,
			expectedErr: false,
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
			err := ur.CreateOne(context.Background(), tt.venue)
			assert.Nil(t, err)

			// Test the GetAll functionality
			_, err = ur.GetAll(context.Background())
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
