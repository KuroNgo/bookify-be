package unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor"
	venue_repository "bookify/internal/repository/venue/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestUpdateOneVenue(t *testing.T) {
	client, database := infrastructor.SetupTestDatabase(t)
	defer infrastructor.TearDownTestDatabase(client, t)
	// Function to clear the event collection before each test case
	clearEventTypeCollection := func() {
		err := database.Collection("venue").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear venue collection: %v", err)
		}
	}

	clearEventTypeCollection()
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

	ur := venue_repository.NewVenueRepository(database, "venue")
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

		err = ur.UpdateOne(context.Background(), mockVenueUpdate)
		assert.Nil(t, err)
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
