package event_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByID
// @Summary Get event by ID
// @Description Get details of an event by its ID
// @Tags Events
// @Accept json
// @Produce json
// @Param id query string true "Event ID"
// @Router /api/v1/events/get/id [get]
func (e *EventController) GetByID(ctx *gin.Context) {
	id := ctx.Query("id")

	data, err := e.EventUseCase.GetByID(ctx, id)
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
