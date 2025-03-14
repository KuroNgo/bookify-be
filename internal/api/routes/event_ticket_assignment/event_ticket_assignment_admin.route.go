package event_ticket_assignment_route

import (
	event_ticket_assignment_controller "bookify/internal/api/controller/event_ticket_assignment"
	"bookify/internal/config"
	"bookify/internal/domain"
	event_repository "bookify/internal/repository/event/repository"
	eventticketrepository "bookify/internal/repository/event_ticket/repository"
	event_ticket_assignment_repository "bookify/internal/repository/event_ticket_assignment/repository"
	userrepository "bookify/internal/repository/user/repository"
	event_ticket_assignment_usecase "bookify/internal/usecase/event_ticket_assignment/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminEventTicketAssignmentRouter(env *config.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	evta := event_ticket_assignment_repository.NewEventTicketAssignmentRepository(db, domain.CollectionEmployeeTicketAssigment)
	ur := userrepository.NewUserRepository(db, domain.CollectionUser)
	evt := eventticketrepository.NewEventTicketRepository(db, domain.CollectionEventTicket)
	ev := event_repository.NewEventRepository(db, domain.CollectionEvent)

	eventType := &event_ticket_assignment_controller.EventTicketAssignmentController{
		EventTicketAssignmentUseCase: event_ticket_assignment_usecase.NewEventTicketAssignmentUseCase(env, timeout, evta, ev, evt, ur),
		Database:                     env,
	}

	router := group.Group("/event-ticket-assignments")
	router.POST("/create", eventType.CreateOne)
	router.PUT("/update", eventType.UpdateOne)
	router.POST("/delete", eventType.DeleteOne)
	router.GET("/statistics-revenue/event_id", eventType.StatisticsRevenueByEventID)
}
