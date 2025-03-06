package unit

import (
	"bookify/internal/domain"
	"bookify/internal/infrastructor/mongodb"
	activity_log_repository "bookify/internal/repository/activity_log/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestCreateOneActivityLog(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	// Function to clear the activity log collection before each test case
	clearActivityLogCollection := func() {
		err := database.Collection("activity_logs").Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to clear activity_logs collection: %v", err)
		}
	}

	// Mock data
	mockActivityLog := &domain.ActivityLog{
		ID:           primitive.NewObjectID(),
		ClientIP:     "192.168.1.1",
		UserID:       primitive.NewObjectID(),
		Level:        1,
		Method:       "POST",
		StatusCode:   200,
		BodySize:     512,
		Path:         "/api/v1/login",
		Latency:      "100ms",
		ActivityTime: time.Now(),
		ExpireAt:     time.Now().Add(24 * time.Hour),
	}

	// Define test cases
	tests := []struct {
		name        string
		input       *domain.ActivityLog
		expectedErr bool
		description string
	}{
		{
			name:        "success_create_activity_log",
			input:       mockActivityLog,
			expectedErr: false,
			description: "Should successfully create an activity log",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearActivityLogCollection() // Clear the collection before each test

			repo := activity_log_repository.NewActivityLogRepository(database, "activity_logs")
			err := repo.CreateOne(context.Background(), tt.input)

			if tt.expectedErr {
				assert.Error(t, err, tt.description)
			} else {
				assert.Nil(t, err, tt.description)
			}
		})
	}
}
