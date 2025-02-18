package event_controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByStartTime
// @Summary Get event by start time
// @Description Get event starting from a specific date
// @Tags Events
// @Accept json
// @Produce json
// @Param startDate query string true "Start date in YYYY-MM-DD format"
// @Router /api/v1/events/get/start-time [get]
func (e *EventController) GetByStartTime(ctx *gin.Context) {
	startDate := ctx.Query("startTime")
	if startDate == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing startDate parameter"})
		return
	}

	//Xử lý nếu startDate là JSON
	var startDateData map[string]string
	if err := json.Unmarshal([]byte(startDate), &startDateData); err == nil {
		startDate = startDateData["startTime"]
	}

	// Gọi UseCase để xử lý logic
	data, err := e.EventUseCase.GetByStartTime(ctx, startDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to get event: " + err.Error(),
		})
		return
	}

	// Trả về kết quả
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Events retrieved successfully",
		"data":    data,
	})
}
