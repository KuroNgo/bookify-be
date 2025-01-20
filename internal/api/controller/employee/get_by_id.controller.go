package employee_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByID godoc
// @Summary Get an employee by ID
// @Description Retrieve an employee record by its ID
// @Tags Employees
// @Accept json
// @Produce json
// @Param id query string true "Employee ID"
// @Router /api/v1/employees/get/id [get]
func (e EmployeeController) GetByID(ctx *gin.Context) {
	id := ctx.Query("id")
	data, err := e.EmployeeUseCase.GetByID(ctx, id)
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
