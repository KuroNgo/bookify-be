package event_ticket_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAll godoc
// @Summary Get all event tickets
// @Description Retrieves a list of all event tickets
// @Tags Event Tickets
// @Accept json
// @Produce json
// @Router /api/v1/event-tickets/get/all [get]
func (e *EventTicketController) GetAll(ctx *gin.Context) {
	data, err := e.EventTicketUseCase.GetAll(ctx)
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
