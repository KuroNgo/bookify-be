package event_controller

import (
	"bookify/pkg/shared/constants"
	"github.com/gin-gonic/gin"
	"net/http"
)

// DeleteOne
// @Summary Delete an event
// @Description Delete an event by its ID
// @Tags Events
// @Accept json
// @Produce json
// @Param id query string true "Event ID"
// @Router /api/v1/events/delete-one [delete]
func (e *EventController) DeleteOne(ctx *gin.Context) {
	_, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": constants.MsgAPIUnauthorized,
		})
		return
	}

	id := ctx.Query("id")
	err := e.EventUseCase.DeleteOne(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to get user data: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
