package event_type

import (
	event_type_controller "bookify/internal/api/controller/event_type"
	"bookify/internal/config"
	"bookify/internal/domain"
	event_type_repository "bookify/internal/repository/event_type/repository"
	event_type_usecase "bookify/internal/usecase/event_type/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func EventTypeRouter(env *config.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ev := event_type_repository.NewEventTypeRepository(db, domain.CollectionEventType)

	event_type := &event_type_controller.EventController{
		EventTypeUseCase: event_type_usecase.NewEventTypeUseCase(env, timeout, ev),
		Database:         env,
	}

	router := group.Group("/event-types")
	router.GET("/get/id", event_type.GetByID)
	router.GET("/get/all", event_type.GetAll)
}
