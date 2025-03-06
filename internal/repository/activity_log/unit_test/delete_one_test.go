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

func TestDeleteOneActivityLog(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	repo := activity_log_repository.NewActivityLogRepository(database, "activity_logs")

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

	err := repo.CreateOne(context.Background(), mockActivityLog)
	assert.Nil(t, err, "Should successfully create an activity log before deleting")

	t.Run("success_delete_activity_log", func(t *testing.T) {
		err := repo.DeleteOne(context.Background(), mockActivityLog.ID)
		assert.Nil(t, err, "Should successfully delete an existing activity log")
	})
}
