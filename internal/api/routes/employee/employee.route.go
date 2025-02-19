package employee_route

import (
	employee_controller "bookify/internal/api/controller/employee"
	"bookify/internal/config"
	"bookify/internal/domain"
	employee_repository "bookify/internal/repository/employee/repository"
	organizationrepository "bookify/internal/repository/organization/repository"
	user_repository "bookify/internal/repository/user/repository"
	employee_usecase "bookify/internal/usecase/employee/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func EmployeeRouter(env *config.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	em := employee_repository.NewEmployeeRepository(db, domain.CollectionEmployee)
	ur := user_repository.NewUserRepository(db, domain.CollectionUser)
	or := organizationrepository.NewOrganizationRepository(db, domain.CollectionOrganization)

	employee := &employee_controller.EmployeeController{
		EmployeeUseCase: employee_usecase.NewEmployeeUseCase(env, timeout, em, ur, or),
		Database:        env,
	}

	router := group.Group("/employees")
	router.GET("/get/id", employee.GetByID)
	router.GET("/get/all", employee.GetAll)
}
