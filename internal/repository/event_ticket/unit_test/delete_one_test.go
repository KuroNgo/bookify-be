package unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor"
	eventticketrepository "bookify/internal/repository/event_ticket/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestDeleteOneEventTicket(t *testing.T) {
	client, database := infrastructor.SetupTestDatabase(t)
	defer infrastructor.TearDownTestDatabase(client, t)

	// Hàm dọn dẹp collection trước mỗi test
	clearEventTicketCollection := func() {
		err := database.Collection("event_ticket").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear event_ticket collection: %v", err)
		}
	}

	clearEventTicketCollection()

	// Khởi tạo repository
	ur := eventticketrepository.NewEventTicketRepository(database, "event_ticket")

	// Tạo vé mẫu trong database
	eventTicket := domain.EventTicket{
		ID:       primitive.NewObjectID(),
		EventID:  primitive.NewObjectID(),
		Price:    100.0,
		Quantity: 10,
		Status:   "available",
	}
	_, err := database.Collection("event_ticket").InsertOne(context.Background(), eventTicket)
	assert.NoError(t, err)

	tests := []struct {
		name      string
		id        primitive.ObjectID
		expectErr bool
	}{
		{
			name:      "success",
			id:        eventTicket.ID,
			expectErr: false,
		},
		{
			name:      "error_invalid_id",
			id:        primitive.NewObjectID(), // ID không tồn tại
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ur.DeleteOne(context.Background(), tt.id)

			if tt.expectErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
		})
	}
}
