package event_wishlist_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// DeleteOne godoc
// @Summary Delete an event wishlist
// @Description Deletes an event wishlist by ID for the current user
// @Tags Event Wishlists
// @Accept json
// @Produce json
// @Param id query string true "Event Wishlist ID"
// @Router /api/v1/event-wishlists/delete [delete]
func (e *EventWishlistController) DeleteOne(ctx *gin.Context) {
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "You are not login!",
		})
		return
	}

	id := ctx.Query("id")
	err := e.EventWishlistUseCase.DeleteOne(ctx, id, fmt.Sprintf("%s", currentUser))
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
