package event_type_unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor"
	event_type_repository "bookify/internal/repository/event_type/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestUpdateOneEventType(t *testing.T) {
	client, database := infrastructor.SetupTestDatabase(t)
	defer infrastructor.TearDownTestDatabase(client, t)
	// Function to clear the event collection before each test case
	clearEventTypeCollection := func() {
		err := database.Collection("event_type").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear partner collection: %v", err)
		}
	}

	clearEventTypeCollection()
	mockEventType := domain.EventType{
		ID:   primitive.NewObjectID(),
		Name: "music",
	}

	ur := event_type_repository.NewEventTypeRepository(database, "event_type")
	err := ur.CreateOne(context.Background(), mockEventType)
	assert.Nil(t, err)

	t.Run("success", func(t *testing.T) {
		mockEventTypeUpdate := domain.EventType{
			ID:   mockEventType.ID,
			Name: "food",
		}
		err = ur.UpdateOne(context.Background(), mockEventTypeUpdate)
		assert.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockEventTypeUpdate := domain.EventType{
			ID:   primitive.NilObjectID,
			Name: "music",
		}
		err = ur.UpdateOne(context.Background(), mockEventTypeUpdate)
		assert.Error(t, err)
	})
}
