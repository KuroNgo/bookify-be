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

func TestCreateOneVenue(t *testing.T) {
	client, database := infrastructor.SetupTestDatabase(t)
	defer infrastructor.TearDownTestDatabase(client, t)

	// Function to clear the event collection before each test case
	clearVenueCollection := func() {
		err := database.Collection("venue").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear partner collection: %v", err)
		}
	}

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
	mockVenueNil := &domain.Venue{}

	t.Run("success", func(t *testing.T) {
		clearVenueCollection() // Clear the collection
		ur := venue_repository.NewVenueRepository(database, "venue")
		err := ur.CreateOne(context.Background(), mockVenue)
		assert.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		clearVenueCollection() // Clear the collection
		ur := venue_repository.NewVenueRepository(database, "venue")
		err := ur.CreateOne(context.Background(), mockVenueNil)
		assert.Error(t, err)
	})
}
