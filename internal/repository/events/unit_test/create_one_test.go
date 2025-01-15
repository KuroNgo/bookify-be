package unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor"
	event_repository "bookify/internal/repository/events/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestCreateOneEvent(t *testing.T) {
	client, database := infrastructor.SetupTestDatabase(t)
	defer infrastructor.TearDownTestDatabase(client, t)

	// Function to clear the event collection before each test case
	clearEventCollection := func() {
		err := database.Collection("event").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear employee collection: %v", err)
		}
	}

	mockEvent := &domain.Event{
		EventTypeID:       primitive.NewObjectID(),
		VenueID:           primitive.NewObjectID(),
		OrganizationID:    primitive.NewObjectID(),
		Title:             "title",
		Description:       "description",
		StartTime:         time.Now(),
		EndTime:           time.Now().Add(time.Hour + 2),
		Mode:              "public",
		EstimatedAttendee: 50,
		ActualAttendee:    40,
		TotalExpenditure:  120000000,
	}

	t.Run("success", func(t *testing.T) {
		clearEventCollection() // Clear the collection
		ur := event_repository.NewEventRepository(database, "event")
		err := ur.CreateOne(context.Background(), mockEvent)
		assert.Nil(t, err)
	})
}
