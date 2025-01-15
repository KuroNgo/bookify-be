package event_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAllPagination
// @Summary Get all events with pagination
// @Description Get a paginated list of all events
// @Tags Events
// @Accept json
// @Produce json
// @Param page query string false "Page number, default is 1"
// @Router /api/v1/events/get-all/pagination [get]
func (e *EventController) GetAllPagination(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")

	data, pageOutput, currentPage, err := e.EventUseCase.GetAllPagination(ctx, page)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to get user data: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"data":         data,
		"page":         pageOutput,
		"current_page": currentPage,
	})
}
