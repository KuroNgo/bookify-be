package event_controller

import (
	"bookify/internal/domain"
	"bookify/pkg/shared/constants"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateOneAsync
// @Summary Create a new event
// @Description Add a new event to the system
// @Tags Events
// @Accept json
// @Produce json
// @Param event body domain.EventInput true "Event input data"
// @Router /api/v1/event/create/async [post]
func (e *EventController) CreateOneAsync(ctx *gin.Context) {
	_, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": constants.MsgAPIUnauthorized,
		})
		return
	}

	//  Lấy thông tin từ request
	var eventInput domain.EventInput
	if err := ctx.ShouldBindJSON(&eventInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	err := e.EventUseCase.CreateOneAsync(ctx, &eventInput)
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
