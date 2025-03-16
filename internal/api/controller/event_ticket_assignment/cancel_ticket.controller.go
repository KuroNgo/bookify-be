package event_ticket_assignment_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CancelTickets godoc
// @Summary Cancel an event ticket assignment for cancel tickets
// @Description Cancel an existing event ticket assignment
// @Tags Event Ticket Assignments
// @Accept json
// @Produce json
// @Param id query string true "Event Ticket Assignment ID"
// @Router /api/v1/event-ticket-assignments/update/cancel-ticket [patch]
func (e EventTicketAssignmentController) CancelTickets(ctx *gin.Context) {
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "You are not login!",
		})
		return
	}
	id := ctx.Query("id")

	err := e.EventTicketAssignmentUseCase.CancelTickets(ctx, id, fmt.Sprintf("%s", currentUser))
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
