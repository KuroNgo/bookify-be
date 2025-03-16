package event_ticket_assignment_route

import (
	eventticketassignmentcontroller "bookify/internal/api/controller/event_ticket_assignment"
	"bookify/internal/config"
	"bookify/internal/domain"
	eventrepository "bookify/internal/repository/event/repository"
	eventticketrepository "bookify/internal/repository/event_ticket/repository"
	eventticketassignmentrepository "bookify/internal/repository/event_ticket_assignment/repository"
	userrepository "bookify/internal/repository/user/repository"
	eventticketassignmentusecase "bookify/internal/usecase/event_ticket_assignment/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminEventTicketAssignmentRouter(env *config.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	evta := eventticketassignmentrepository.NewEventTicketAssignmentRepository(db, domain.CollectionEmployeeTicketAssigment)
	ur := userrepository.NewUserRepository(db, domain.CollectionUser)
	evt := eventticketrepository.NewEventTicketRepository(db, domain.CollectionEventTicket)
	ev := eventrepository.NewEventRepository(db, domain.CollectionEvent)

	eventType := &eventticketassignmentcontroller.EventTicketAssignmentController{
		EventTicketAssignmentUseCase: eventticketassignmentusecase.NewEventTicketAssignmentUseCase(env, timeout, evta, ev, evt, ur),
		Database:                     env,
	}

	router := group.Group("/event-ticket-assignments")
	router.POST("/create", eventType.CreateOne)
	router.PUT("/update", eventType.UpdateOne)
	router.PATCH("/update/cancel-ticket", eventType.CancelTickets)
	router.POST("/delete", eventType.DeleteOne)
	router.GET("/statistics-revenue/event_id", eventType.StatisticsRevenueByEventID)
}
