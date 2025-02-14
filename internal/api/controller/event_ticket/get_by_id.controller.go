package event_ticket_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByID godoc
// @Summary Get event tickets by ID
// @Description Retrieves the details of an event tickets by its ID
// @Tags Event Tickets
// @Accept json
// @Produce json
// @Param id query string true "Event Ticket ID"
// @Router /api/v1/event-tickets/get/id [get]
func (e *EventTicketController) GetByID(ctx *gin.Context) {
	id := ctx.Query("id")
	data, err := e.EventTicketUseCase.GetByID(ctx, id)
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
