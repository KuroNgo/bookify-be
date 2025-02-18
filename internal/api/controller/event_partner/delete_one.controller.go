package event_partner_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// DeleteOne godoc
// @Summary Delete an event partner
// @Description Deletes an event partner by ID for the current user
// @Tags Event Partners
// @Accept json
// @Produce json
// @Param id query string true "Event Partner ID"
// @Router /api/v1/event-partners/delete [delete]
func (e *EventPartnerController) DeleteOne(ctx *gin.Context) {
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "You are not login!",
		})
		return
	}

	id := ctx.Query("id")
	err := e.EventPartnerUseCase.DeleteOne(ctx, id, fmt.Sprintf("%s", currentUser))
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
