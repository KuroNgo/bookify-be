package event_ticket_assignment_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByID godoc
// @Summary Get event ticket assignment by ID
// @Description Retrieves the details of an event ticket assignment by its ID
// @Tags Event Ticket Assignments
// @Accept json
// @Produce json
// @Param id query string true "Event Ticket Assignment ID"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Router /api/v1/event-ticket-assignments/get/id [get]
func (e EventTicketAssignmentController) GetByID(ctx *gin.Context) {
	id := ctx.Query("id")
	data, err := e.EventTicketAssignmentUseCase.GetByID(ctx, id)
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
