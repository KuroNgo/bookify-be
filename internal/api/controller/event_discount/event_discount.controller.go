package event_discount_controller

import (
	"bookify/internal/config"
	eventdiscountusecase "bookify/internal/usecase/event_discount/usecase"
	cronjob "bookify/pkg/shared/schedules"
	"github.com/gin-gonic/gin"
)

type EventDiscountController struct {
	Database             *config.Database
	CronJob              *cronjob.CronScheduler
	EventDiscountUseCase eventdiscountusecase.IEventDiscountUseCase
}

type IEventTypeUseCase interface {
	GetByID(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	CreateOne(ctx *gin.Context)
	UpdateOne(ctx *gin.Context)
	DeleteOne(ctx *gin.Context)
}

func NewEvent(Database *config.Database, EventDiscountUseCase eventdiscountusecase.IEventDiscountUseCase) IEventTypeUseCase {
	return &EventDiscountController{Database: Database, EventDiscountUseCase: EventDiscountUseCase}
}
