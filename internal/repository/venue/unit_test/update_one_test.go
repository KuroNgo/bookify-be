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

func TestUpdateOneVenue(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	// Function to clear the venue collection before each test case
	clearVenueCollection := func() {
		err := database.Collection("venue").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear venue collection: %v", err)
		}
	}

	clearVenueCollection()
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

	ur := venuerepository.NewVenueRepository(database, "venue")
	err := ur.CreateOne(context.Background(), mockVenue)
	assert.Nil(t, err)

	t.Run("success", func(t *testing.T) {
		mockVenueUpdate := &domain.Venue{
			ID:          mockVenue.ID,
			Capacity:    120,
			AddressLine: "123 Main Street",
			City:        "New York",
			State:       "NY",
			Country:     "USA",
			PostalCode:  "10001",
			OnlineFlat:  false,
		}

		// Update the venue
		err = ur.UpdateOne(context.Background(), mockVenueUpdate)
		assert.Nil(t, err)

		// Fetch the updated venue and assert the changes
		updatedVenue, err := ur.GetByID(context.Background(), mockVenue.ID)
		assert.Nil(t, err)
		assert.Equal(t, mockVenueUpdate.Capacity, updatedVenue.Capacity)
		assert.Equal(t, mockVenueUpdate.AddressLine, updatedVenue.AddressLine)
	})

	t.Run("error", func(t *testing.T) {
		mockVenueUpdateNil := &domain.Venue{
			ID:          primitive.NilObjectID,
			Capacity:    120,
			AddressLine: "123 Main Street",
			City:        "New York",
			State:       "NY",
			Country:     "USA",
			PostalCode:  "10001",
			OnlineFlat:  false,
		}
		err = ur.UpdateOne(context.Background(), mockVenueUpdateNil)
		assert.Error(t, err)
	})
}
