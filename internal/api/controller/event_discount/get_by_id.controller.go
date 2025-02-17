package event_discount_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByID godoc
// @Summary Get event discount by ID
// @Description Retrieves the details of an event discount by its ID
// @Tags Event Discounts
// @Accept json
// @Produce json
// @Param id query string true "Event Discount ID"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Router /api/v1/event-discount/get/id [get]
func (e *EventDiscountController) GetByID(ctx *gin.Context) {
	id := ctx.Query("id")
	data, err := e.EventDiscountUseCase.GetByID(ctx, id)
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
