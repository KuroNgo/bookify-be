package event_ticket_assignment_controller

import (
	"bookify/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateOne godoc
// @Summary Create a new event ticket assignment
// @Description Creates a new event ticket assignment for the current user
// @Tags Event Ticket Assignments
// @Accept json
// @Produce json
// @Param eventTicketAssignment body domain.EventTicketAssignmentInput true "Event Ticket Assignment Body"
// @Router /api/v1/event-ticket-assignments/create [post]
func (e EventTicketAssignmentController) CreateOne(ctx *gin.Context) {
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

	err := e.EventTicketAssignmentUseCase.CreateOne(ctx, &eventTicketAssignmentInput, fmt.Sprintf("%s", currentUser))
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
