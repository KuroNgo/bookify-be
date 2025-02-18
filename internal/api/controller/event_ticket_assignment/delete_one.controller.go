package event_ticket_assignment_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// DeleteOne godoc
// @Summary Delete an event ticket assignment
// @Description Deletes an event ticket assignment by ID for the current user
// @Tags Event Ticket Assignments
// @Accept json
// @Produce json
// @Param id query string true "Event Ticket Assignment ID"
// @Router /api/v1/event-ticket-assignment/delete [delete]
func (e EventTicketAssignmentController) DeleteOne(ctx *gin.Context) {
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "You are not login!",
		})
		return
	}

	id := ctx.Query("id")
	err := e.EventTicketAssignmentUseCase.DeleteOne(ctx, id, fmt.Sprintf("%s", currentUser))
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
