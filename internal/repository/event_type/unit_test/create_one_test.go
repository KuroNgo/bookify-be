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

func TestCreateOneEventType(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	clearEventTypeCollection := func() {
		err := database.Collection("event_type").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear event_type collection: %v", err)
		}
	}

	clearEventTypeCollection()

	tests := []struct {
		name      string
		input     domain.EventType
		expectErr bool
	}{
		{
			name:      "success",
			input:     domain.EventType{ID: primitive.NewObjectID(), Name: "music"},
			expectErr: false,
		},
		{
			name:      "error",
			input:     domain.EventType{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := event_type_repository.NewEventTypeRepository(database, "event_type")
			err := ur.CreateOne(context.Background(), tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
