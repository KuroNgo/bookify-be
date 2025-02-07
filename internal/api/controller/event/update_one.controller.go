package event_controller

import (
	"bookify/internal/domain"
	"bookify/pkg/shared/constants"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UpdateOne
// @Summary Update an event
// @Description Update details of an existing event
// @Tags Events
// @Accept json
// @Produce json
// @Param id query string true "Event ID"
// @Param event body domain.EventInput true "Event input data"
// @Router /api/v1/event/update [put]
func (e *EventController) UpdateOne(ctx *gin.Context) {
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
	id := ctx.Query("id")

	err := e.EventUseCase.UpdateOne(ctx, id, &eventInput)
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
