package partner_controller

import (
	"bookify/internal/config"
	partner_usecase "bookify/internal/usecase/partner/usecase"
	"github.com/gin-gonic/gin"
)

type PartnerController struct {
	Database       *config.Database
	PartnerUseCase partner_usecase.IPartnerUseCase
}

type IPartnerUseCase interface {
	GetByID(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	CreateOne(ctx *gin.Context)
	UpdateOne(ctx *gin.Context)
	DeleteOne(ctx *gin.Context)
}

func NewEvent(Database *config.Database, PartnerUseCase partner_usecase.IPartnerUseCase) IPartnerUseCase {
	return &PartnerController{Database: Database, PartnerUseCase: PartnerUseCase}
}
