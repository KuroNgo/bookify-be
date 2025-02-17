package event_wishlist_controller

import (
	"bookify/internal/config"
	event_wishlist_usecase "bookify/internal/usecase/event_wishlist/usecase"
	"github.com/gin-gonic/gin"
)

type EventWishlistController struct {
	Database             *config.Database
	EventWishlistUseCase event_wishlist_usecase.IEventWishlistUseCase
}

type IEventWishlistUseCase interface {
	GetByID(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	CreateOne(ctx *gin.Context)
	UpdateOne(ctx *gin.Context)
	DeleteOne(ctx *gin.Context)
}

func NewEventWishlist(Database *config.Database, EventWishlistUseCase event_wishlist_usecase.IEventWishlistUseCase) IEventWishlistUseCase {
	return &EventWishlistController{Database: Database, EventWishlistUseCase: EventWishlistUseCase}
}
