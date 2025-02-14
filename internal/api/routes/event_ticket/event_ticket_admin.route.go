package event_ticket_route

import (
	event_ticket_controller "bookify/internal/api/controller/event_ticket"
	"bookify/internal/config"
	"bookify/internal/domain"
	event_repository "bookify/internal/repository/event/repository"
	eventticketrepository "bookify/internal/repository/event_ticket/repository"
	userrepository "bookify/internal/repository/user/repository"
	event_ticket_usecase "bookify/internal/usecase/event_ticket/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminEventTicketRouter(env *config.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	evt := eventticketrepository.NewEventTicketRepository(db, domain.CollectionEventTicket)
	ev := event_repository.NewEventRepository(db, domain.CollectionEvent)
	ur := userrepository.NewUserRepository(db, domain.CollectionUser)

	eventTicket := &event_ticket_controller.EventTicketController{
		EventTicketUseCase: event_ticket_usecase.NewEventTicketUseCase(env, timeout, evt, ev, ur),
		Database:           env,
	}

	router := group.Group("/event-tickets")
	router.POST("/create", eventTicket.CreateOne)
	router.PUT("/update", eventTicket.UpdateOne)
	router.POST("/delete", eventTicket.DeleteOne)
}
