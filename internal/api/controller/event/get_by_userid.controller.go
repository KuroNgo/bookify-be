package event_controller

import (
	"bookify/pkg/shared/constants"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByUserID
// @Summary Get event by UserID
// @Description Get details of an event by its User ID
// @Tags Events
// @Accept json
// @Produce json
// @Router /api/v1/event/get/user_id [get]
func (e *EventController) GetByUserID(ctx *gin.Context) {
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": constants.MsgAPIUnauthorized,
		})
		return
	}

	data, err := e.EventUseCase.GetByUserID(ctx, fmt.Sprintf("%s", currentUser))
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
