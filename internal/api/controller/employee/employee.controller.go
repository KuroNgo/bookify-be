package employee_controller

import (
	"bookify/internal/config"
	employee_usecase "bookify/internal/usecase/employee/usecase"
	"github.com/gin-gonic/gin"
)

type EmployeeController struct {
	Database        *config.Database
	EmployeeUseCase employee_usecase.IEmployeeUseCase
}

type IEmployeeController interface {
	GetByID(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	CreateOne(ctx *gin.Context)
	UpdateOne(ctx *gin.Context)
	DeleteOne(ctx *gin.Context)
}

func NewEmployee(Database *config.Database, EmployeeUseCase employee_usecase.IEmployeeUseCase) IEmployeeController {
	return &EmployeeController{Database: Database, EmployeeUseCase: EmployeeUseCase}
}
