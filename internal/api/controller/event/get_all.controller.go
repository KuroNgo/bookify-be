package event_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAll
// @Summary Get all event
// @Description Get a list of all event
// @Tags Events
// @Accept json
// @Produce json
// @Router /api/v1/event/get/all [get]
func (e *EventController) GetAll(ctx *gin.Context) {
	data, err := e.EventUseCase.GetAll(ctx)
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
