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

func TestCreateOneEventTicket(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	// Hàm dọn dẹp collection trước mỗi test
	clearEventTicketCollection := func() {
		err := database.Collection("event_ticket").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear event_ticket collection: %v", err)
		}
	}

	clearEventTicketCollection()

	tests := []struct {
		name      string
		input     domain.EventTicket
		expectErr bool
	}{
		{
			name: "success",
			input: domain.EventTicket{
				ID:       primitive.NewObjectID(),
				EventID:  primitive.NewObjectID(),
				Price:    100.0,
				Quantity: 10,
				Status:   "available",
			},
			expectErr: false,
		},
		{
			name:      "error_missing_fields",
			input:     domain.EventTicket{}, // Trường hợp thiếu dữ liệu
			expectErr: true,
		},
	}

	// Khởi tạo repository
	ur := eventticketrepository.NewEventTicketRepository(database, "event_ticket")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ur.CreateOne(context.Background(), tt.input)

			if tt.expectErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
		})
	}
}
