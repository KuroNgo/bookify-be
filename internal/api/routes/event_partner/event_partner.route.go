package event_partner_route

import (
	event_partner_controller "bookify/internal/api/controller/event_partner"
	"bookify/internal/config"
	"bookify/internal/domain"
	event_partner_repository "bookify/internal/repository/event_partner/repository"
	userrepository "bookify/internal/repository/user/repository"
	event_partner_usecase "bookify/internal/usecase/event_partner/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func EventPartnerRouter(env *config.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ev := event_partner_repository.NewEventTypeRepository(db, domain.CollectionEventPartner)
	ur := userrepository.NewUserRepository(db, domain.CollectionUser)

	eventPartner := &event_partner_controller.EventPartnerController{
		EventPartnerUseCase: event_partner_usecase.NewEventPartnerUseCase(env, timeout, ev, ur),
		Database:            env,
	}

	router := group.Group("/event-partners")
	router.GET("/get/id", eventPartner.GetByID)
	router.GET("/get/all", eventPartner.GetAll)
}
