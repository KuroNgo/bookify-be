package event_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByStartTimePagination
// @Summary Get event by start time with pagination
// @Description Get paginated event starting from a specific date
// @Tags Events
// @Accept json
// @Produce json
// @Param startDate query string true "Start date in YYYY-MM-DD format"
// @Param page query string false "Page number, default is 1"
// @Router /api/v1/events/get/start-time/pagination [get]
func (e *EventController) GetByStartTimePagination(ctx *gin.Context) {
	startDate := ctx.Query("startDate")
	page := ctx.DefaultQuery("page", "1")

	data, err := e.EventUseCase.GetByStartTimePagination(ctx, startDate, page)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to get user data: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"data":         data,
		"page":         data.Page,
		"current_page": data.CurrentPage,
	})
}
