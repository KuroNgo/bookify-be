package event_employee_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByID godoc
// @Summary Get event employee by ID
// @Description Retrieves the details of an event employee by its ID
// @Tags Event Employees
// @Accept json
// @Produce json
// @Param id query string true "Event Employee ID"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Router /api/v1/event-employees/get/id [get]
func (e EventEmployeeController) GetByID(ctx *gin.Context) {
	id := ctx.Query("id")
	data, err := e.EventEmployeeUseCase.GetByID(ctx, id)
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
