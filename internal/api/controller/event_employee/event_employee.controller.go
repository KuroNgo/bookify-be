package event_employee_controller

import (
	"bookify/internal/config"
	event_employee_usecase "bookify/internal/usecase/event_employee/usecase"
	"github.com/gin-gonic/gin"
)

type EventEmployeeController struct {
	Database             *config.Database
	EventEmployeeUseCase event_employee_usecase.IEventEmployeeUseCase
}

type IEventEmployeeUseCase interface {
	GetByID(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	CreateOne(ctx *gin.Context)
	UpdateOne(ctx *gin.Context)
	DeleteOne(ctx *gin.Context)
}

func NewEvent(Database *config.Database, EventEmployeeUseCase event_employee_usecase.IEventEmployeeUseCase) IEventEmployeeUseCase {
	return &EventEmployeeController{Database: Database, EventEmployeeUseCase: EventEmployeeUseCase}
}
