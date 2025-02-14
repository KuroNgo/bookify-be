package unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor"
	eventticketrepository "bookify/internal/repository/event_ticket/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestUpdateOneEventTicket(t *testing.T) {
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

	// Tạo một vé mẫu trong database
	existingTicket := domain.EventTicket{
		ID:       primitive.NewObjectID(),
		EventID:  primitive.NewObjectID(),
		Price:    100.0,
		Quantity: 10,
		Status:   "available",
	}
	_, err := database.Collection("event_ticket").InsertOne(context.Background(), existingTicket)
	assert.NoError(t, err)

	tests := []struct {
		name      string
		input     domain.EventTicket
		expectErr bool
	}{
		{
			name: "success_update_price_and_status",
			input: domain.EventTicket{
				ID:      existingTicket.ID,
				EventID: existingTicket.EventID,
				Price:   150.0,
				Status:  "sold_out",
			},
			expectErr: false,
		},
		{
			name: "error_invalid_id",
			input: domain.EventTicket{
				ID:     primitive.NewObjectID(), // ID không tồn tại
				Price:  200.0,
				Status: "available",
			},
			expectErr: true,
		},
		{
			name: "error_empty_id",
			input: domain.EventTicket{
				Price:  180.0,
				Status: "expired",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ur.UpdateOne(context.Background(), tt.input)

			if tt.expectErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")

				// Kiểm tra dữ liệu đã cập nhật trong DB
				var updatedTicket domain.EventTicket
				err := database.Collection("event_ticket").FindOne(context.Background(), bson.M{"_id": tt.input.ID}).Decode(&updatedTicket)
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Price, updatedTicket.Price)
				assert.Equal(t, tt.input.Status, updatedTicket.Status)
			}
		})
	}
}
