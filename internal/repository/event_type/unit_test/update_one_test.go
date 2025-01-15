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

func TestUpdateOneEventType(t *testing.T) {
	client, database := infrastructor.SetupTestDatabase(t)
	defer infrastructor.TearDownTestDatabase(client, t)

	mockEventType := &domain.EventType{
		ID:            primitive.NewObjectID(),
		EventTypeName: "music",
	}

	ur := event_type_repository.NewEventTypeRepository(database, "event_type")
	err := ur.CreateOne(context.Background(), mockEventType)
	assert.Nil(t, err)

	t.Run("success", func(t *testing.T) {
		mockEventTypeUpdate := &domain.EventType{
			ID:            mockEventType.ID,
			EventTypeName: "music",
		}

		err = ur.UpdateOne(context.Background(), mockEventTypeUpdate)
		assert.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockEventTypeUpdate := &domain.EventType{
			ID:            primitive.NilObjectID,
			EventTypeName: "music",
		}
		err = ur.UpdateOne(context.Background(), mockEventTypeUpdate)
		assert.Error(t, err)
	})
}
