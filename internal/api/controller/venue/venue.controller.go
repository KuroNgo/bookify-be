package venue_controller

import (
	"bookify/internal/config"
	venueusecase "bookify/internal/usecase/venue/usecase"
	"github.com/gin-gonic/gin"
)

type VenueController struct {
	Database     *config.Database
	VenueUseCase venueusecase.IVenueUseCase
}

type IVenueUseCase interface {
	GetByID(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	CreateOne(ctx *gin.Context)
	UpdateOne(ctx *gin.Context)
	DeleteOne(ctx *gin.Context)
}

func NewVenue(Database *config.Database, VenueUseCase venueusecase.IVenueUseCase) IVenueUseCase {
	return &VenueController{Database: Database, VenueUseCase: VenueUseCase}
}
