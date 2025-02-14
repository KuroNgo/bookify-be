package event_ticket_controller

import (
	"bookify/internal/config"
	event_ticket_usecase "bookify/internal/usecase/event_ticket/usecase"
	"github.com/gin-gonic/gin"
)

type EventTicketController struct {
	Database           *config.Database
	EventTicketUseCase event_ticket_usecase.IEventTicketUseCase
}

type IEventTicketUseCase interface {
	GetByID(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	CreateOne(ctx *gin.Context)
	UpdateOne(ctx *gin.Context)
	DeleteOne(ctx *gin.Context)
}

func NewEvent(Database *config.Database, EventTicketUseCase event_ticket_usecase.IEventTicketUseCase) IEventTicketUseCase {
	return &EventTicketController{Database: Database, EventTicketUseCase: EventTicketUseCase}
}
