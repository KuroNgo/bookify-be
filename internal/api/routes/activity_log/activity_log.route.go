package activity_log_route

import (
	activity_log_controller "bookify/internal/api/controller/activity_log"
	"bookify/internal/config"
	"bookify/internal/domain"
	activitylogrepository "bookify/internal/repository/activity_log/repository"
	userrepository "bookify/internal/repository/user/repository"
	activitylogusecase "bookify/internal/usecase/activity_log/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ActivityRoute(env *config.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ac := activitylogrepository.NewActivityLogRepository(db, domain.CollectionActivityLog)
	users := userrepository.NewUserRepository(db, domain.CollectionUser)

	activity := &activity_log_controller.ActivityController{
		ActivityUseCase: activitylogusecase.NewActivityUseCase(env, timeout, ac, users),
		Database:        env,
	}

	router := group.Group("/activity-logs")
	router.GET("/get/all", activity.GetAll)
	router.GET("/get/id", activity.GetByID)
	router.GET("/get/level", activity.GetByLevel)
	router.GET("/get/user_id", activity.GetByUserID)
}
