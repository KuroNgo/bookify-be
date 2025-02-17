package event_employee_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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
