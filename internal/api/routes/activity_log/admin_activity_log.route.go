package activity_log_route

import (
	activity_log_controller "bookify/internal/api/controller/activity_log"
	"bookify/internal/config"
	"bookify/internal/domain"
	activity_log_repository "bookify/internal/repository/activity_log/repository"
	userrepository "bookify/internal/repository/user/repository"
	activity_log_usecase "bookify/internal/usecase/activity_log/usecase"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func Activity(env *config.Database, client *mongo.Client, timeout time.Duration, db *mongo.Database) *activity_log_controller.ActivityController {
	ac := activity_log_repository.NewActivityLogRepository(db, domain.CollectionActivityLog)
	users := userrepository.NewUserRepository(db, domain.CollectionUser)

	activity := &activity_log_controller.ActivityController{
		ActivityUseCase: activity_log_usecase.NewActivityUseCase(env, timeout, ac, users),
		Database:        env,
	}

	return activity
}
