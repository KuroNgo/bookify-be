package integration

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	"bookify/internal/infrastructor/mongodb"
	venue_repository "bookify/internal/repository/venue/repository"
	venue_usecase "bookify/internal/usecase/venue/usecase"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestVenueUseCase_GetAll(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	databaseConfig := config.Database{}
	venueRepo := venue_repository.NewVenueRepository(database, "venue")
	venueUC := venue_usecase.NewVenueUseCase(&databaseConfig, time.Second*5, venueRepo, nil)

	// Helper: Clear the venue collection before each test_e2e case
	clearVenueCollection := func() {
		err := database.Collection("venue").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear venue collection: %v", err)
		}
	}

	t.Run("Successfully retrieve venue", func(t *testing.T) {
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
		err := venueRepo.CreateOne(context.Background(), mockVenue)
		assert.NoError(t, err)

		_, err = venueUC.GetAll(context.Background())
		assert.NoError(t, err)
	})
}
