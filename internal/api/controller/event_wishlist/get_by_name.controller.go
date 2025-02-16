package event_wishlist_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByName godoc
// @Summary Get event type by Name
// @Description Retrieves the details of an event type by its Name
// @Tags Event Types
// @Accept json
// @Produce json
// @Param name query string true "Event Type Name"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Router /api/v1/event-types/get/name [get]
func (e EventController) GetByName(ctx *gin.Context) {
	id := ctx.Query("name")
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
