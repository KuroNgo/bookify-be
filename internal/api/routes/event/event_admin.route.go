package event

import (
	event_controller "bookify/internal/api/controller/event"
	"bookify/internal/config"
	"bookify/internal/domain"
	event_repository "bookify/internal/repository/event/repository"
	eventtyperepository "bookify/internal/repository/event_type/repository"
	organizationrepository "bookify/internal/repository/organization/repository"
	user_repository "bookify/internal/repository/user/repository"
	venue_repository "bookify/internal/repository/venue/repository"
	"bookify/internal/usecase/event/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminEventsRouter(env *config.Database, timeout time.Duration, db *mongo.Database, client *mongo.Client, group *gin.RouterGroup) {
	ev := event_repository.NewEventRepository(db, domain.CollectionEvent)
	or := organizationrepository.NewOrganizationRepository(db, domain.CollectionOrganization)
	evt := eventtyperepository.NewEventTypeRepository(db, domain.CollectionEventType)
	ve := venue_repository.NewVenueRepository(db, domain.CollectionVenue)
	ur := user_repository.NewUserRepository(db, domain.CollectionUser)

	event := &event_controller.EventController{
		EventUseCase: event_usecase.NewEventUseCase(env, timeout, ev, or, evt, ve, ur, client),
		Database:     env,
	}

	router := group.Group("/events")
	router.POST("/create", event.CreateOne)
	router.POST("/create/async", event.CreateOneAsync)
	router.PUT("/update", event.UpdateOne)
	router.PATCH("/update/image", event.UpdateImage)
	router.POST("/delete", event.DeleteOne)
}
