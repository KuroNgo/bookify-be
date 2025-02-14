package venue_controller

import (
	"bookify/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UpdateOne godoc
// @Summary Update venue details
// @Description Update the details of a specific venue by providing venue data and the current user
// @Tags Venues
// @Accept  json
// @Produce  json
// @Param venueInput body domain.VenueInput true "Venue data to update" // Body chứa thông tin venue
// @Security ApiKeyAuth
// @Router /api/v1/venues/update [put]
func (v VenueController) UpdateOne(ctx *gin.Context) {
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "You are not login!",
		})
		return
	}

	//  Lấy thông tin từ request
	var venueInput domain.VenueInput
	if err := ctx.ShouldBindJSON(&venueInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}
	id := ctx.Query("id")

	err := v.VenueUseCase.UpdateOne(ctx, id, &venueInput, fmt.Sprintf("%s", currentUser))
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
