package event_type_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByID godoc
// @Summary Get event type by ID
// @Description Retrieves the details of an event type by its ID
// @Tags Event Types
// @Accept json
// @Produce json
// @Param id query string true "Event Type ID"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Router /api/v1/event-types/get/id [get]
func (e EventController) GetByID(ctx *gin.Context) {
	id := ctx.Query("id")
	data, err := e.EventTypeUseCase.GetByID(ctx, id)
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
