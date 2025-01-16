package partner_route

import (
	partner_controller "bookify/internal/api/controller/partner"
	"bookify/internal/config"
	"bookify/internal/domain"
	partner_repository "bookify/internal/repository/partner/repository"
	user_repository "bookify/internal/repository/user/repository"
	partner_usecase "bookify/internal/usecase/partner/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func PartnerRouter(env *config.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	pn := partner_repository.NewPartnerRepository(db, domain.CollectionPartner)
	ur := user_repository.NewUserRepository(db, domain.CollectionUser)

	partner := &partner_controller.PartnerController{
		PartnerUseCase: partner_usecase.NewPartnerUseCase(env, timeout, pn, ur),
		Database:       env,
	}

	router := group.Group("/partners")
	router.GET("/get/id", partner.GetByID)
	router.GET("/get/all", partner.GetAll)
}
