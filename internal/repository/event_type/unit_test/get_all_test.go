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

func TestFindAllEventType(t *testing.T) {
	client, database := infrastructor.SetupTestDatabase(t)
	defer infrastructor.TearDownTestDatabase(client, t)

	// Function to clear the event collection before each test case
	clearEventTypeCollection := func() {
		err := database.Collection("event_type").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear event_type collection: %v", err)
		}
	}

	clearEventTypeCollection()

	mockEventType := domain.EventType{
		ID:   primitive.NewObjectID(),
		Name: "music",
	}

	ur := event_type_repository.NewEventTypeRepository(database, "event_type")

	// Define test cases
	tests := []struct {
		name         string
		setupFunc    func()
		expectedErr  bool
		expectedSize int
	}{
		{
			name: "success_with_one_event",
			setupFunc: func() {
				err := ur.CreateOne(context.Background(), mockEventType)
				assert.Nil(t, err)
			},
			expectedErr:  false,
			expectedSize: 1,
		},
		{
			name: "success_no_event",
			setupFunc: func() {
				// Không tạo event nào
			},
			expectedErr:  false,
			expectedSize: 0,
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear collection before each test
			clearEventTypeCollection()

			// Run setup function
			tt.setupFunc()

			// Call the function under test
			events, err := ur.GetAll(context.Background())

			// Assert the results
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Len(t, events, tt.expectedSize, "Unexpected number of event")
			}
		})
	}
}
