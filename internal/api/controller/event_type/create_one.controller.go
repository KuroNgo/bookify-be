package event_type_controller

import (
	"bookify/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateOne godoc
// @Summary Create a new event type
// @Description Creates a new event type for the current user
// @Tags Event Types
// @Accept json
// @Produce json
// @Param eventType body domain.EventTypeInput true "Event Type Body"
// @Router /api/v1/event-types/create [post]
func (e EventController) CreateOne(ctx *gin.Context) {
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "You are not login!",
		})
		return
	}

	//  Lấy thông tin từ request
	var userInput domain.EventTypeInput
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	err := e.EventTypeUseCase.CreateOne(ctx, &userInput, fmt.Sprintf("%s", currentUser))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	// Trả về phản hồi thành công
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
