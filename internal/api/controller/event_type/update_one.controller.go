package event_type_controller

import (
	"bookify/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UpdateOne godoc
// @Summary Update an event type
// @Description Updates an existing event type
// @Tags Event Types
// @Accept json
// @Produce json
// @Param eventType body domain.EventTypeInput true "Event Type Body"
// @Router /api/v1/event-types/update [put]
func (e EventController) UpdateOne(ctx *gin.Context) {
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "You are not login!",
		})
		return
	}

	//  Lấy thông tin từ request
	var eventTypeInput domain.EventTypeInput
	if err := ctx.ShouldBindJSON(&eventTypeInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}
	id := ctx.Query("id")

	err := e.EventTypeUseCase.UpdateOne(ctx, id, &eventTypeInput, fmt.Sprintf("%s", currentUser))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
