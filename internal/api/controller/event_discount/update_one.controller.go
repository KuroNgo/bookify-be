package event_discount_controller

import (
	"bookify/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UpdateOne godoc
// @Summary Update an event discount
// @Description Updates an existing event discount
// @Tags Event Discounts
// @Accept json
// @Produce json
// @Param eventDiscount body domain.EventDiscount true "Event Discount Body"
// @Router /api/v1/event-discount/update [put]
func (e *EventDiscountController) UpdateOne(ctx *gin.Context) {
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "You are not login!",
		})
		return
	}

	//  Lấy thông tin từ request
	var eventDiscountInput domain.EventDiscountInput
	if err := ctx.ShouldBindJSON(&eventDiscountInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}
	id := ctx.Query("id")

	err := e.EventDiscountUseCase.UpdateOne(ctx, id, &eventDiscountInput, fmt.Sprintf("%s", currentUser))
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
