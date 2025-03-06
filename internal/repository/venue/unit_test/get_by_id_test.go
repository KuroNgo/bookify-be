package unit_test

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor/mongodb"
	venuerepository "bookify/internal/repository/venue/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestFindByIDVenue(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	// Function to clear the venue collection before each test_e2e case
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
		Country:     "USA",
		EventMode:   "online",
		LinkAttend:  "",
		FromAttend:  "",
	}
	ur := venuerepository.NewVenueRepository(database, "venue")
	err := ur.CreateOne(context.Background(), mockVenue)
	assert.Nil(t, err)

	t.Run("success", func(t *testing.T) {
		_, err = ur.GetByID(context.Background(), mockVenue.ID)
		assert.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		_, err = ur.GetByID(context.Background(), primitive.NilObjectID)
		assert.Error(t, err)
	})
}
