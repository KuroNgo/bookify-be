package unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor/mongodb"
	eventticketrepository "bookify/internal/repository/event_ticket/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestGetByIDEventTicket(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	// Hàm dọn dẹp collection trước mỗi test_e2e
	clearEventTicketCollection := func() {
		err := database.Collection("event_ticket").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear event_ticket collection: %v", err)
		}
	}

	clearEventTicketCollection()

	// Khởi tạo repository
	ur := eventticketrepository.NewEventTicketRepository(database, "event_ticket")

	// Tạo một vé mẫu trong database
	existingTicket := domain.EventTicket{
		ID:       primitive.NewObjectID(),
		EventID:  primitive.NewObjectID(),
		Price:    200.0,
		Quantity: 5,
		Status:   "available",
	}
	_, err := database.Collection("event_ticket").InsertOne(context.Background(), existingTicket)
	assert.NoError(t, err)

	tests := []struct {
		name      string
		inputID   primitive.ObjectID
		expectErr bool
		expectNil bool
	}{
		{
			name:      "success_get_existing_ticket",
			inputID:   existingTicket.ID,
			expectErr: false,
			expectNil: false,
		},
		{
			name:      "error_invalid_id",
			inputID:   primitive.NilObjectID,
			expectErr: true,
			expectNil: false,
		},
		{
			name:      "error_ticket_not_found",
			inputID:   primitive.NewObjectID(), // ID không tồn tại
			expectErr: true,
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ur.GetByID(context.Background(), tt.inputID)

			if tt.expectErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")

				if tt.expectNil {
					assert.Equal(t, domain.EventTicket{}, result, "Expected empty result but got data")
				} else {
					assert.Equal(t, existingTicket.ID, result.ID, "Expected matching ID but got different one")
					assert.Equal(t, existingTicket.Price, result.Price)
					assert.Equal(t, existingTicket.Quantity, result.Quantity)
					assert.Equal(t, existingTicket.Status, result.Status)
				}
			}
		})
	}
}
