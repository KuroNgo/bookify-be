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

func TestGetByUserIDActivityLog(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	repo := activity_log_repository.NewActivityLogRepository(database, "activity_logs")

	userID := primitive.NewObjectID()

	mockActivityLog1 := &domain.ActivityLog{
		ID:           primitive.NewObjectID(),
		ClientIP:     "192.168.1.3",
		UserID:       userID,
		Level:        2,
		Method:       "DELETE",
		StatusCode:   204,
		BodySize:     128,
		Path:         "/api/v1/logout",
		Latency:      "30ms",
		ActivityTime: time.Now(),
		ExpireAt:     time.Now().Add(24 * time.Hour),
	}

	mockActivityLog2 := &domain.ActivityLog{
		ID:           primitive.NewObjectID(),
		ClientIP:     "192.168.1.4",
		UserID:       userID,
		Level:        3,
		Method:       "PATCH",
		StatusCode:   200,
		BodySize:     300,
		Path:         "/api/v1/update-profile",
		Latency:      "80ms",
		ActivityTime: time.Now(),
		ExpireAt:     time.Now().Add(24 * time.Hour),
	}

	err := repo.CreateOne(context.Background(), mockActivityLog1)
	assert.Nil(t, err, "Should successfully create first user activity log")

	err = repo.CreateOne(context.Background(), mockActivityLog2)
	assert.Nil(t, err, "Should successfully create second user activity log")

	t.Run("success_get_activity_logs_by_user_id", func(t *testing.T) {
		logs, err := repo.GetByUserID(context.Background(), userID)
		assert.Nil(t, err, "Should successfully retrieve activity logs by user ID")
		assert.Len(t, logs, 2, "Should return 2 logs for the given user ID")
	})

	t.Run("error_get_activity_logs_by_nonexistent_user_id", func(t *testing.T) {
		nonexistentUserID := primitive.NewObjectID()
		logs, err := repo.GetByUserID(context.Background(), nonexistentUserID)
		assert.Nil(t, err, "Should not return an error when no logs are found")
		assert.Empty(t, logs, "Should return an empty slice when no logs match the user ID")
	})
}
