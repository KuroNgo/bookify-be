package partner_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAll godoc
// @Summary Get a list of all partners
// @Description Retrieve a list of all partners from the system
// @Tags partners
// @Accept json
// @Produce json
// @Router /api/v1/partners/get/all [get]
func (p PartnerController) GetAll(ctx *gin.Context) {
	data, err := p.PartnerUseCase.GetAll(ctx)
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
