package event_ticket_assignment_controller

import (
	"bookify/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UpdateOne godoc
// @Summary Update an event ticket assignment
// @Description Updates an existing event ticket assignment
// @Tags Event Ticket Assignments
// @Accept json
// @Produce json
// @Param eventTicketAssignment body domain.EventTicketAssignmentInput true "Event Ticket Assignment Body"
// @Router /api/v1/event-ticket-assignments/update [put]
func (e EventTicketAssignmentController) UpdateOne(ctx *gin.Context) {
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "You are not login!",
		})
		return
	}

	//  Lấy thông tin từ request
	var eventTicketAssignmentInput domain.EventTicketAssignmentInput
	if err := ctx.ShouldBindJSON(&eventTicketAssignmentInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}
	id := ctx.Query("id")

	err := e.EventTicketAssignmentUseCase.UpdateOne(ctx, id, &eventTicketAssignmentInput, fmt.Sprintf("%s", currentUser))
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
