package event_controller

import (
	"bookify/pkg/shared/constants"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByUserIDAndStartTime
// @Summary Get events by User ID and Start Time
// @Description Retrieve event details for a specific user based on User ID and Start Time
// @Tags Events
// @Accept json
// @Produce json
// @Param start_time query string true "Start Time (ISO 8601 format)"
// @Security BearerAuth
// @Router /api/v1/events/get/user_id/start_time [get]
func (e *EventController) GetByUserIDAndStartTime(ctx *gin.Context) {
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": constants.MsgAPIUnauthorized,
		})
		return
	}
	startTime := ctx.Query("start_time")

	data, err := e.EventUseCase.GetByOrganizationIDAndStartTime(ctx, fmt.Sprintf("%s", currentUser), startTime)
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
