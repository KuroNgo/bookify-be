package organization_controller

import (
	"bookify/internal/config"
	organization_usecase "bookify/internal/usecase/organization/usecase"
	"github.com/gin-gonic/gin"
)

type OrganizationController struct {
	Database            *config.Database
	OrganizationUseCase organization_usecase.IOrganizationUseCase
}

type IOrganizationController interface {
	GetByID(ctx *gin.Context)
	GetByUserID(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	CreateOne(ctx *gin.Context)
	UpdateOne(ctx *gin.Context)
	DeleteOne(ctx *gin.Context)
}

func NewOrganization(Database *config.Database, OrganizationUseCase organization_usecase.IOrganizationUseCase) IOrganizationController {
	return &OrganizationController{Database: Database, OrganizationUseCase: OrganizationUseCase}
}
