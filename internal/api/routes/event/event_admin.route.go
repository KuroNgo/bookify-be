package event

import (
	event_controller "bookify/internal/api/controller/event"
	"bookify/internal/config"
	"bookify/internal/domain"
	event_repository "bookify/internal/repository/events/repository"
	"bookify/internal/usecase/event/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminEventsRouter(env *config.Database, timeout time.Duration, db *mongo.Database, client *mongo.Client, group *gin.RouterGroup) {
	ev := event_repository.NewEventRepository(db, domain.CollectionEvent)

	event := &event_controller.EventController{
		EventUseCase: usecase.NewEventUseCase(env, timeout, ev, client),
		Database:     env,
	}

	router := group.Group("/events")
	router.POST("/create", event.CreateOne)
	router.PUT("/update", event.UpdateOne)
	router.POST("/delete", event.DeleteOne)
}
