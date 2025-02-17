package event_discount_controller

import (
	"bookify/internal/config"
	event_discount_usecase "bookify/internal/usecase/event_discount/usecase"
	"github.com/gin-gonic/gin"
)

type EventDiscountController struct {
	Database             *config.Database
	EventDiscountUseCase event_discount_usecase.IEventDiscountUseCase
}

type IEventTypeUseCase interface {
	GetByID(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	CreateOne(ctx *gin.Context)
	UpdateOne(ctx *gin.Context)
	DeleteOne(ctx *gin.Context)
}

func NewEvent(Database *config.Database, EventDiscountUseCase event_discount_usecase.IEventDiscountUseCase) IEventTypeUseCase {
	return &EventDiscountController{Database: Database, EventDiscountUseCase: EventDiscountUseCase}
}
