package event_partner_controller

import (
	"bookify/internal/config"
	event_partner_usecase "bookify/internal/usecase/event_partner/usecase"
	"github.com/gin-gonic/gin"
)

type EventPartnerController struct {
	Database            *config.Database
	EventPartnerUseCase event_partner_usecase.IEventPartnerUseCase
}

type IEventPartnerUseCase interface {
	GetByID(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	CreateOne(ctx *gin.Context)
	UpdateOne(ctx *gin.Context)
	DeleteOne(ctx *gin.Context)
}

func NewEvent(Database *config.Database, EventPartnerUseCase event_partner_usecase.IEventPartnerUseCase) IEventPartnerUseCase {
	return &EventPartnerController{Database: Database, EventPartnerUseCase: EventPartnerUseCase}
}
