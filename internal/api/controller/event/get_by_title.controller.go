package event_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByTitle
// @Summary Get event by title
// @Description Get details of an event by title
// @Tags Events
// @Accept json
// @Produce json
// @Router /api/v1/events/get/title [get]
func (e *EventController) GetByTitle(ctx *gin.Context) {
	title := ctx.Query("title")

	data, err := e.EventUseCase.GetByTitle(ctx, title)
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
