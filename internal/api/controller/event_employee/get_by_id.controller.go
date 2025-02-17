package event_employee_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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
