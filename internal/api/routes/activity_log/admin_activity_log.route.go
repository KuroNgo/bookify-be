package activity_log_route

import (
	activitylogcontroller "bookify/internal/api/controller/activity_log"
	"bookify/internal/config"
	"bookify/internal/domain"
	activitylogrepository "bookify/internal/repository/activity_log/repository"
	userrepository "bookify/internal/repository/user/repository"
	activitylogusecase "bookify/internal/usecase/activity_log/usecase"
	userusecase "bookify/internal/usecase/user/usecase"
	cronjob "bookify/pkg/shared/schedules"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func Activity(env *config.Database, cr *cronjob.CronScheduler, client *mongo.Client, timeout time.Duration, db *mongo.Database) *activitylogcontroller.ActivityController {
	ac := activitylogrepository.NewActivityLogRepository(db, domain.CollectionActivityLog)
	users := userrepository.NewUserRepository(db, domain.CollectionUser)

	activity := &activitylogcontroller.ActivityController{
		ActivityUseCase: activitylogusecase.NewActivityUseCase(env, cr, timeout, ac, users),
		UserUseCase:     userusecase.NewUserUseCase(env, timeout, users, client),
		Database:        env,
	}

	return activity
}
