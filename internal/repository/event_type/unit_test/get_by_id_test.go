package unit_test

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

	mockEventType := &domain.EventType{
		ID:            primitive.NewObjectID(),
		EventTypeName: "music",
	}

	ur := event_type_repository.NewEventTypeRepository(database, "event_type")
	err := ur.CreateOne(context.Background(), mockEventType)
	assert.Nil(t, err)

	t.Run("success", func(t *testing.T) {
		_, err = ur.GetByID(context.Background(), mockEventType.ID)
		assert.Nil(t, err)
	})

	t.Run("success", func(t *testing.T) {
		_, err = ur.GetByID(context.Background(), primitive.NilObjectID)
		assert.Nil(t, err)
	})
}
