package unit

import (
	"bookify/internal/infrastructor/mongodb"
	activity_log_repository "bookify/internal/repository/activity_log/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAllActivityLogs(t *testing.T) {
	client, database := mongodb.SetupTestDatabase(t)
	defer mongodb.TearDownTestDatabase(client, t)

	repo := activity_log_repository.NewActivityLogRepository(database, "activity_logs")

	t.Run("success_get_all_activity_logs", func(t *testing.T) {
		logs, err := repo.GetAll(context.Background())
		assert.Nil(t, err, "Should successfully retrieve all activity logs")
		assert.NotNil(t, logs, "Logs should not be nil")
	})
}
