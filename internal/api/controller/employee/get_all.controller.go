package employee_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAll godoc
// @Summary Get all employees
// @Description Retrieve all employee records
// @Tags Employees
// @Accept json
// @Produce json
// @Router /api/v1/employees/get/all [get]
func (e EmployeeController) GetAll(ctx *gin.Context) {
	data, err := e.EmployeeUseCase.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to get user data: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}
