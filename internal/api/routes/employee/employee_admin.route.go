package employee_route

import (
	employee_controller "bookify/internal/api/controller/employee"
	"bookify/internal/config"
	"bookify/internal/domain"
	employee_repository "bookify/internal/repository/employee/repository"
	user_repository "bookify/internal/repository/user/repository"
	employee_usecase "bookify/internal/usecase/employee/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminEmployeeRouter(env *config.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	em := employee_repository.NewEmployeeRepository(db, domain.CollectionEmployee)
	ur := user_repository.NewUserRepository(db, domain.CollectionUser)

	employee := &employee_controller.EmployeeController{
		EmployeeUseCase: employee_usecase.NewEmployeeUseCase(env, timeout, em, ur),
		Database:        env,
	}

	router := group.Group("/employees")
	router.POST("/create", employee.CreateOne)
	router.PUT("/update", employee.UpdateOne)
	router.DELETE("/delete", employee.DeleteOne)
}
