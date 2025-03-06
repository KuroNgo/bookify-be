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

func TestUpdateOneEventType(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	// Function to clear the event collection before each test_e2e case
	clearEventTypeCollection := func() {
		err := database.Collection("event_type").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear event type collection: %v", err)
		}
	}

	clearEventTypeCollection()

	// Mock data
	mockEventType := domain.EventType{
		ID:   primitive.NewObjectID(),
		Name: "music",
	}

	ur := event_type_repository.NewEventTypeRepository(database, "event_type")
	err := ur.CreateOne(context.Background(), mockEventType)
	assert.Nil(t, err)

	// Define test_e2e cases
	tests := []struct {
		name        string
		input       domain.EventType
		expectedErr bool
		description string
	}{
		{
			name: "success_update_existing_event_type",
			input: domain.EventType{
				ID:   mockEventType.ID,
				Name: "food",
			},
			expectedErr: false,
			description: "Should successfully update an existing event type",
		},
		{
			name: "error_update_invalid_id",
			input: domain.EventType{
				ID:   primitive.NilObjectID,
				Name: "music",
			},
			expectedErr: true,
			description: "Should return an error when updating with invalid ID",
		},
	}

	// Execute test_e2e cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ur.UpdateOne(context.Background(), tt.input)

			if tt.expectedErr {
				assert.Error(t, err, tt.description)
			} else {
				assert.Nil(t, err, tt.description)
			}
		})
	}
}
