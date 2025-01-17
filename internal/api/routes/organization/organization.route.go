package organization_route

import (
	organizationcontroller "bookify/internal/api/controller/organization"
	"bookify/internal/config"
	"bookify/internal/domain"
	organizationrepository "bookify/internal/repository/organization/repository"
	userrepository "bookify/internal/repository/user/repository"
	organizationusecase "bookify/internal/usecase/organization/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func OrganizationRouter(env *config.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	or := organizationrepository.NeOrganizationRepository(db, domain.CollectionOrganization)
	ur := userrepository.NewUserRepository(db, domain.CollectionUser)

	organization := &organizationcontroller.OrganizationController{
		OrganizationUseCase: organizationusecase.NewOrganizationUseCase(env, timeout, or, ur),
		Database:            env,
	}

	router := group.Group("/organization")
	router.GET("/get/id", organization.GetByID)
	router.GET("/get/all", organization.GetAll)
}
