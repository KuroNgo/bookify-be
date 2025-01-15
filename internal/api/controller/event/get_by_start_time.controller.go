package event_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByStartTime
// @Summary Get events by start time
// @Description Get events starting from a specific date
// @Tags Events
// @Accept json
// @Produce json
// @Param startDate query string true "Start date in YYYY-MM-DD format"
// @Router /api/v1/events/get/start-time [get]
func (e *EventController) GetByStartTime(ctx *gin.Context) {
	startDate := ctx.Query("startDate")

	data, err := e.EventUseCase.GetByStartTime(ctx, startDate)
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
