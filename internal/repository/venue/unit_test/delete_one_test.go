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

func TestDeleteOneVenue(t *testing.T) {
	client, database := infrastructor.SetupTestDatabase(t)
	defer infrastructor.TearDownTestDatabase(client, t)

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
		err = ur.DeleteOne(context.Background(), mockVenue.ID)
		assert.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		err = ur.DeleteOne(context.Background(), primitive.NilObjectID)
		assert.Error(t, err)
	})
}
