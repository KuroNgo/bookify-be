package event_wishlist_route

import (
	event_wishlist_controller "bookify/internal/api/controller/event_wishlist"
	"bookify/internal/config"
	"bookify/internal/domain"
	event_wishlist_repository "bookify/internal/repository/event_wishlist/repository"
	userrepository "bookify/internal/repository/user/repository"
	event_wishlist_usecase "bookify/internal/usecase/event_wishlist/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminEventWishlistRouter(env *config.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	evw := event_wishlist_repository.NewEventWishlistRepository(db, domain.CollectionEventWishlist)
	ur := userrepository.NewUserRepository(db, domain.CollectionUser)

	eventWishlist := &event_wishlist_controller.EventWishlistController{
		EventWishlistUseCase: event_wishlist_usecase.NewEventWishlistUseCase(env, timeout, evw, ur),
		Database:             env,
	}

	router := group.Group("/event-wishlists")
	router.POST("/create", eventWishlist.CreateOne)
	router.PUT("/update", eventWishlist.UpdateOne)
	router.POST("/delete", eventWishlist.DeleteOne)
}
