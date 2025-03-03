package event_employee_route

import (
	event_employee_controller "bookify/internal/api/controller/event_employee"
	"bookify/internal/config"
	"bookify/internal/domain"
	employee_repository "bookify/internal/repository/employee/repository"
	event_employee_repository "bookify/internal/repository/event_employee/repository"
	userrepository "bookify/internal/repository/user/repository"
	event_employee_usecase "bookify/internal/usecase/event_employee/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func EventEmployeeRouter(env *config.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	eve := event_employee_repository.NewEventEmployeeRepository(db, domain.CollectionEventEmployee)
	ur := userrepository.NewUserRepository(db, domain.CollectionUser)
	em := employee_repository.NewEmployeeRepository(db, domain.CollectionEmployee)

	eventEmployee := &event_employee_controller.EventEmployeeController{
		EventEmployeeUseCase: event_employee_usecase.NewEventEmployeeUseCase(env, timeout, eve, em, ur),
		Database:             env,
	}

	router := group.Group("/event-employees")
	router.GET("/get/id", eventEmployee.GetByID)
	router.GET("/get/all", eventEmployee.GetAll)
}
