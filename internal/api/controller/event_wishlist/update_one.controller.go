package event_wishlist_controller

import (
	"bookify/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UpdateOne godoc
// @Summary Update an event wishlists
// @Description Updates an existing event wishlists
// @Tags Event Wishlists
// @Accept json
// @Produce json
// @Param eventWishlist body domain.EventWishlistInput true "Event Wishlist Body"
// @Router /api/v1/event-wishlist/update [put]
func (e *EventWishlistController) UpdateOne(ctx *gin.Context) {
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "You are not login!",
		})
		return
	}

	//  Lấy thông tin từ request
	var eventWishlistInput domain.EventWishlistInput
	if err := ctx.ShouldBindJSON(&eventWishlistInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}
	id := ctx.Query("id")

	err := e.EventWishlistUseCase.UpdateOne(ctx, id, &eventWishlistInput, fmt.Sprintf("%s", currentUser))
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
