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

func TestGetByLevelActivityLog(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	repo := activity_log_repository.NewActivityLogRepository(database, "activity_logs")

	mockActivityLog1 := &domain.ActivityLog{
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

	mockActivityLog2 := &domain.ActivityLog{
		ID:           primitive.NewObjectID(),
		ClientIP:     "192.168.1.2",
		UserID:       primitive.NewObjectID(),
		Level:        1,
		Method:       "GET",
		StatusCode:   200,
		BodySize:     256,
		Path:         "/api/v1/profile",
		Latency:      "50ms",
		ActivityTime: time.Now(),
		ExpireAt:     time.Now().Add(24 * time.Hour),
	}

	err := repo.CreateOne(context.Background(), mockActivityLog1)
	assert.Nil(t, err, "Should successfully create first activity log")

	err = repo.CreateOne(context.Background(), mockActivityLog2)
	assert.Nil(t, err, "Should successfully create second activity log")

	t.Run("success_get_activity_logs_by_level", func(t *testing.T) {
		_, err := repo.GetByLevel(context.Background(), "1")
		assert.Nil(t, err, "Should successfully retrieve activity logs by level")
		assert.NoError(t, err, "Should successfully retrieve activity logs by level")
	})

	t.Run("error_get_activity_logs_by_nonexistent_level", func(t *testing.T) {
		logs, err := repo.GetByLevel(context.Background(), "999")
		assert.Nil(t, err, "Should not return an error when no logs are found")
		assert.Empty(t, logs, "Should return an empty slice when no logs match the level")
	})
}
