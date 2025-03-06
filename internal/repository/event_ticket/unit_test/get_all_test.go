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

func TestGetAllEventTickets(t *testing.T) {
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

	// Chuẩn bị dữ liệu test_e2e
	tickets := []domain.EventTicket{
		{
			ID:       primitive.NewObjectID(),
			EventID:  primitive.NewObjectID(),
			Price:    100.0,
			Quantity: 10,
			Status:   "available",
		},
		{
			ID:       primitive.NewObjectID(),
			EventID:  primitive.NewObjectID(),
			Price:    200.0,
			Quantity: 5,
			Status:   "sold_out",
		},
	}

	// Chèn dữ liệu vào database
	var docs []interface{}
	for _, ticket := range tickets {
		docs = append(docs, ticket)
	}
	_, err := database.Collection("event_ticket").InsertMany(context.Background(), docs)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		expectErr  bool
		expectSize int
	}{
		{
			name:       "success_get_all_tickets",
			expectErr:  false,
			expectSize: len(tickets),
		},
		{
			name:       "success_get_empty_list",
			expectErr:  false,
			expectSize: 0, // Nếu collection rỗng
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectSize == 0 {
				// Dọn sạch collection nếu test_e2e case yêu cầu danh sách rỗng
				err := database.Collection("event_ticket").Drop(context.Background())
				assert.NoError(t, err)
			}

			results, err := ur.GetAll(context.Background())

			if tt.expectErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
				assert.Len(t, results, tt.expectSize, "Unexpected number of event tickets")
			}
		})
	}
}
