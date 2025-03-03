package event_employee_route

import (
	event_employee_controller "bookify/internal/api/controller/event_employee"
	"bookify/internal/config"
	"bookify/internal/domain"
	event_employee_repository "bookify/internal/repository/event_employee/repository"
	userrepository "bookify/internal/repository/user/repository"
	event_employee_usecase "bookify/internal/usecase/event_employee/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminEventEmployeeRouter(env *config.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	eve := event_employee_repository.NewEventEmployeeRepository(db, domain.CollectionEventEmployee)
	ur := userrepository.NewUserRepository(db, domain.CollectionUser)

	eventEmployee := &event_employee_controller.EventEmployeeController{
		EventEmployeeUseCase: event_employee_usecase.NewEventEmployeeUseCase(env, timeout, eve, ur),
		Database:             env,
	}

	router := group.Group("/event-employees")
	router.POST("/create", eventEmployee.CreateOne)
	router.PUT("/update", eventEmployee.UpdateOne)
	router.POST("/delete", eventEmployee.DeleteOne)
}
