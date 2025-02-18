package event_ticket_assignment_controller

import (
	"bookify/internal/config"
	event_ticket_assignment_usecase "bookify/internal/usecase/event_ticket_assignment/usecase"
	"github.com/gin-gonic/gin"
)

type EventTicketAssignmentController struct {
	Database                     *config.Database
	EventTicketAssignmentUseCase event_ticket_assignment_usecase.IEventTicketAssignmentUseCase
}

type IEventTicketAssignmentUseCase interface {
	GetByID(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	CreateOne(ctx *gin.Context)
	UpdateOne(ctx *gin.Context)
	DeleteOne(ctx *gin.Context)
}

func NewEvent(Database *config.Database, EventTicketAssignmentUseCase event_ticket_assignment_usecase.IEventTicketAssignmentUseCase) IEventTicketAssignmentUseCase {
	return &EventTicketAssignmentController{Database: Database, EventTicketAssignmentUseCase: EventTicketAssignmentUseCase}
}
