package unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor"
	event_type_repository "bookify/internal/repository/event_type/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestCreateOneEventType(t *testing.T) {
	client, database := infrastructor.SetupTestDatabase(t)
	defer infrastructor.TearDownTestDatabase(client, t)

	// Function to clear the event collection before each test case
	clearEventCollection := func() {
		err := database.Collection("event_type").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear event type collection: %v", err)
		}
	}

	mockEventType := &domain.EventType{
		ID:            primitive.NewObjectID(),
		EventTypeName: "music",
	}

	mockEventTypeNil := &domain.EventType{
		EventTypeName: "",
	}

	t.Run("success", func(t *testing.T) {
		clearEventCollection() // Clear the collection
		ur := event_type_repository.NewEventTypeRepository(database, "event_type")
		err := ur.CreateOne(context.Background(), mockEventType)
		assert.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		clearEventCollection() // Clear the collection
		ur := event_type_repository.NewEventTypeRepository(database, "event_type")
		err := ur.CreateOne(context.Background(), mockEventTypeNil)
		assert.Error(t, err)
	})
}
