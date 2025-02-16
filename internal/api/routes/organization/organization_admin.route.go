package organization_route

import (
	organization_controller "bookify/internal/api/controller/organization"
	"bookify/internal/config"
	"bookify/internal/domain"
	organization_repository "bookify/internal/repository/organization/repository"
	user_repository "bookify/internal/repository/user/repository"
	organization_usecase "bookify/internal/usecase/organization/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminOrganizationRouter(env *config.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	or := organization_repository.NewOrganizationRepository(db, domain.CollectionOrganization)
	ur := user_repository.NewUserRepository(db, domain.CollectionUser)

	organization := &organization_controller.OrganizationController{
		OrganizationUseCase: organization_usecase.NewOrganizationUseCase(env, timeout, or, ur),
		Database:            env,
	}

	router := group.Group("/organizations")
	router.POST("/create", organization.CreateOne)
	router.PUT("/update", organization.UpdateOne)
	router.POST("/delete", organization.DeleteOne)
}
