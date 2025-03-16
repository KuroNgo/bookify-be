package event_ticket_assignment_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// StatisticsRevenueByEventID godoc
// @Summary Revenue statistics by Event ID
// @Description This API returns the total revenue of an event based on its Event ID.
// @Tags Event Ticket Assignments
// @Accept json
// @Produce json
// @Param event_id query string true "ID of the event to retrieve revenue statistics"
// @Router /api/v1/event-ticket-assignments/statistics-revenue/event_id [get]
func (e EventTicketAssignmentController) StatisticsRevenueByEventID(ctx *gin.Context) {
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "You are not login!",
		})
		return
	}

	user, ok := currentUser.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Invalid user type",
		})
		return
	}
	eventId := ctx.Query("event_id")

	data, err := e.EventTicketAssignmentUseCase.StatisticsRevenueByEventID(ctx, eventId, user)
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
