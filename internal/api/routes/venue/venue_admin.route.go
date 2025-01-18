package venue_route

import (
	venue_controller "bookify/internal/api/controller/venue"
	"bookify/internal/config"
	"bookify/internal/domain"
	user_repository "bookify/internal/repository/user/repository"
	venue_repository "bookify/internal/repository/venue/repository"
	venue_usecase "bookify/internal/usecase/venue/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminVenueRouter(env *config.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ve := venue_repository.NewVenueRepository(db, domain.CollectionVenue)
	ur := user_repository.NewUserRepository(db, domain.CollectionUser)

	venue := &venue_controller.VenueController{
		VenueUseCase: venue_usecase.NewVenueUseCase(env, timeout, ve, ur),
		Database:     env,
	}

	router := group.Group("/venues")
	router.POST("/create", venue.CreateOne)
	router.PUT("/update", venue.UpdateOne)
	router.DELETE("/delete", venue.DeleteOne)
}
