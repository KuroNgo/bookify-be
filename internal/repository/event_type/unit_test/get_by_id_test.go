package event_type_unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor/mongodb"
	event_type_repository "bookify/internal/repository/event_type/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestFindByIDEventType(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	// Function to clear the event collection before each test_e2e case
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

	// Define test_e2e cases
	tests := []struct {
		name        string
		inputID     primitive.ObjectID
		expectedErr bool
	}{
		{
			name:        "success",
			inputID:     mockEventType.ID,
			expectedErr: false,
		},
		{
			name:        "error_invalid_id",
			inputID:     primitive.NilObjectID,
			expectedErr: true,
		},
	}

	// Execute test_e2e cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ur.GetByID(context.Background(), tt.inputID)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
