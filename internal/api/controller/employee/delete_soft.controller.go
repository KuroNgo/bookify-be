package employee_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// DeleteSoftOne godoc
// @Summary Delete an employee
// @Description Delete an employee record by ID
// @Tags Employees
// @Accept json
// @Produce json
// @Param id query string true "Employee ID"
// @Router /api/v1/employees/delete-soft [patch]
func (e EmployeeController) DeleteSoftOne(ctx *gin.Context) {
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "You are not login!",
		})
		return
	}

	id := ctx.Query("id")
	err := e.EmployeeUseCase.DeleteSoft(ctx, id, fmt.Sprintf("%s", currentUser))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
