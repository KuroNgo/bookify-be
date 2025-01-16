package partner_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByID godoc
// @Summary Get a partner by ID
// @Description Retrieve a partner from the system using the partner's ID
// @Tags partners
// @Accept json
// @Produce json
// @Param id query string true "Partner ID"
// @Router /api/v1/partners/get/id [get]
func (p PartnerController) GetByID(ctx *gin.Context) {
	id := ctx.Query("id")
	data, err := p.PartnerUseCase.GetByID(ctx, id)
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
