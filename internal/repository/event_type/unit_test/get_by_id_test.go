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

func TestFindByIDEventType(t *testing.T) {
	client, database := infrastructor.SetupTestDatabase(t)
	defer infrastructor.TearDownTestDatabase(client, t)

	// Function to clear the event collection before each test case
	clearEventCollection := func() {
		err := database.Collection("event_type").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear event type collection: %v", err)
		}
	}
	clearEventCollection()
	mockEventType := domain.EventType{
		ID:   primitive.NewObjectID(),
		Name: "event",
	}

	ur := event_type_repository.NewEventTypeRepository(database, "event_type")
	err := ur.CreateOne(context.Background(), mockEventType)
	assert.Nil(t, err)

	t.Run("success", func(t *testing.T) {
		_, err = ur.GetByID(context.Background(), mockEventType.ID)
		assert.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		_, err = ur.GetByID(context.Background(), primitive.NilObjectID)
		assert.Error(t, err)
	})
}
