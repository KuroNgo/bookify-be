package event_employee_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAll godoc
// @Summary Get all event employee
// @Description Retrieves a list of all event employee
// @Tags Event Employees
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Router /api/v1/event-employees/get/all [get]
func (e EventEmployeeController) GetAll(ctx *gin.Context) {
	data, err := e.EventEmployeeUseCase.GetAll(ctx)
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
