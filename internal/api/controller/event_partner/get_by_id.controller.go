package event_partner_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByID godoc
// @Summary Get event partner by ID
// @Description Retrieves the details of an event partner by its ID
// @Tags Event Partners
// @Accept json
// @Produce json
// @Param id query string true "Event Partner ID"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Router /api/v1/event-partners/get/id [get]
func (e *EventPartnerController) GetByID(ctx *gin.Context) {
	id := ctx.Query("id")
	data, err := e.EventPartnerUseCase.GetByID(ctx, id)
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
