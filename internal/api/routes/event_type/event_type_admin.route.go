package event_type

import (
	eventtypecontroller "bookify/internal/api/controller/event_type"
	"bookify/internal/config"
	"bookify/internal/domain"
	eventtyperepository "bookify/internal/repository/event_type/repository"
	userrepository "bookify/internal/repository/user/repository"
	eventtypeusecase "bookify/internal/usecase/event_type/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminEventTypeRouter(env *config.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ev := eventtyperepository.NewEventTypeRepository(db, domain.CollectionEventType)
	ur := userrepository.NewUserRepository(db, domain.CollectionUser)

	eventType := &eventtypecontroller.EventController{
		EventTypeUseCase: eventtypeusecase.NewEventTypeUseCase(env, timeout, ev, ur),
		Database:         env,
	}

	router := group.Group("/event-types")
	router.POST("/create", eventType.CreateOne)
	router.PUT("/update", eventType.UpdateOne)
	router.POST("/delete", eventType.DeleteOne)
}
