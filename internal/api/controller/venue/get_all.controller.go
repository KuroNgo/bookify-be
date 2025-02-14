package venue_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAll godoc
// @Summary Get all venues
// @Description Retrieve a list of all venues
// @Tags Venues
// @Accept  json
// @Produce  json
// @Router /api/v1/venues/get/all [get]
func (v VenueController) GetAll(ctx *gin.Context) {
	data, err := v.VenueUseCase.GetAll(ctx)
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
