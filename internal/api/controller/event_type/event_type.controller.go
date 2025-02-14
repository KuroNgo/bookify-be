package event_type_controller

import (
	"bookify/internal/config"
	event_type_usecase "bookify/internal/usecase/event_type/usecase"
	"github.com/gin-gonic/gin"
)

type EventController struct {
	Database         *config.Database
	EventTypeUseCase event_type_usecase.IEventTypeUseCase
}

type IEventTypeUseCase interface {
	GetByID(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	CreateOne(ctx *gin.Context)
	UpdateOne(ctx *gin.Context)
	DeleteOne(ctx *gin.Context)
}

func NewEvent(Database *config.Database, EventTypeUseCase event_type_usecase.IEventTypeUseCase) IEventTypeUseCase {
	return &EventController{Database: Database, EventTypeUseCase: EventTypeUseCase}
}
