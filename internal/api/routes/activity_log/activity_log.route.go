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
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ActivityRoute(env *config.Database, cr *cronjob.CronScheduler, client *mongo.Client, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ac := activitylogrepository.NewActivityLogRepository(db, domain.CollectionActivityLog)
	users := userrepository.NewUserRepository(db, domain.CollectionUser)

	activity := &activitylogcontroller.ActivityController{
		ActivityUseCase: activitylogusecase.NewActivityUseCase(env, cr, timeout, ac, users),
		UserUseCase:     userusecase.NewUserUseCase(env, timeout, users, client),
		Database:        env,
	}

	router := group.Group("/activity-logs")
	router.GET("/get/all", activity.GetAll)
	router.GET("/get/id", activity.GetByID)
	router.GET("/get/level", activity.GetByLevel)
	router.GET("/get/user_id", activity.GetByUserID)
}
