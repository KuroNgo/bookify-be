package event_wishlist_controller

import (
	"bookify/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateOne godoc
// @Summary Create a new event wishlist
// @Description Creates a new event wishlist for the current user
// @Tags Event Wishlists
// @Accept json
// @Produce json
// @Param eventWishlist body domain.EventWishlist true "Event Wishlist Body"
// @Router /api/v1/event-wishlists/create [post]
func (e *EventWishlistController) CreateOne(ctx *gin.Context) {
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

	err := e.EventWishlistUseCase.CreateOne(ctx, &eventWishlistInput, fmt.Sprintf("%s", currentUser))
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
