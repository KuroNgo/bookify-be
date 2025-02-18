package event_employee_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// DeleteOne godoc
// @Summary Delete an event employee
// @Description Deletes an event employee by ID for the current user
// @Tags Event Employees
// @Accept json
// @Produce json
// @Param id query string true "Event Type ID"
// @Router /api/v1/event-employees/delete [delete]
func (e EventEmployeeController) DeleteOne(ctx *gin.Context) {
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "You are not login!",
		})
		return
	}

	id := ctx.Query("id")
	err := e.EventEmployeeUseCase.DeleteOne(ctx, id, fmt.Sprintf("%s", currentUser))
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
