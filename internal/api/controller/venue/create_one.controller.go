package venue_controller

import (
	"bookify/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateOne docs
// @Summary Create a new venue
// @Description Create a new venue using the provided input data
// @Tags Venues
// @Accept json
// @Produce json
// @Param venueInput body domain.VenueInput true "Venue input data"
// @Router /api/v1/venues/create [post]
func (v VenueController) CreateOne(ctx *gin.Context) {
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

	err := v.VenueUseCase.CreateOne(ctx, &venueInput, fmt.Sprintf("%s", currentUser))
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
