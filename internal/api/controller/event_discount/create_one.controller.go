package event_discount_controller

import (
	"bookify/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateOne godoc
// @Summary Create a new event discount
// @Description Creates a new event discount for the current user
// @Tags Event Discounts
// @Accept json
// @Produce json
// @Param discount body domain.EventDiscountInput true "Event discount Body"
// @Router /api/v1/event-discount/create [post]
func (e *EventDiscountController) CreateOne(ctx *gin.Context) {
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

	err := e.EventDiscountUseCase.CreateOne(ctx, &eventDiscountInput, fmt.Sprintf("%s", currentUser))
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
