package venue_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByID godoc
// @Summary Get venue by ID
// @Description Retrieve a venue by its ID
// @Tags Venues
// @Accept  json
// @Produce  json
// @Param id query string true "Venue ID"
// @Router /api/v1/venues/get/id [get]
func (v VenueController) GetByID(ctx *gin.Context) {
	id := ctx.Query("id")
	data, err := v.VenueUseCase.GetByID(ctx, id)
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
