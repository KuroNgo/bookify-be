package unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor/mongodb"
	eventwishlistrepository "bookify/internal/repository/event_wishlist/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestGetByUserID(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t) // Thiết lập database test_e2e
	defer mongodb.TearDownTestDatabase(client, t)    // Dọn dẹp sau khi test_e2e xong

	// Mock data
	mockUserID := primitive.NewObjectID()
	mockWishlist := domain.EventWishlist{
		ID:     primitive.NewObjectID(),
		UserID: mockUserID,
	}

	// Khởi tạo repository
	eventWishlistRepo := eventwishlistrepository.NewEventWishlistRepository(database, "event_wishlist")

	// Tạo một document trong database
	err := eventWishlistRepo.CreateOne(context.Background(), mockWishlist)
	assert.Nil(t, err)

	// Define test_e2e cases
	tests := []struct {
		name        string
		userID      primitive.ObjectID
		expectedErr bool
		description string
	}{
		{
			name:        "success_get_existing_wishlist",
			userID:      mockUserID,
			expectedErr: true,
			description: "Should return the event wishlist for a valid user ID",
		},
		{
			name:        "error_get_non_existing_wishlist",
			userID:      primitive.NewObjectID(),
			expectedErr: false, // Vì hàm GetByUserID không trả lỗi khi không tìm thấy tài liệu
			description: "Should return an empty wishlist for a non-existing user ID",
		},
	}

	// Execute test_e2e cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wishlist, err := eventWishlistRepo.GetByUserID(context.Background(), tt.userID)

			// Kiểm tra lỗi có như mong đợi không
			assert.Nil(t, err, tt.description)

			// Nếu tìm thấy dữ liệu thì ID phải khớp
			if tt.userID == mockUserID {
				assert.Equal(t, mockWishlist.ID, wishlist.ID, "EventWishlist ID should match")
			} else {
				// Nếu không tìm thấy thì struct phải rỗng
				assert.Equal(t, domain.EventWishlist{}, wishlist, "Should return an empty EventWishlist")
			}
		})
	}
}
