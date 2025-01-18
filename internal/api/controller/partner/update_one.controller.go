package partner_controller

import (
	"bookify/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UpdateOne godoc
// @Summary Update a partner by ID
// @Description Update the details of a partner using the partner's ID and input data
// @Tags partners
// @Accept json
// @Produce json
// @Param id query string true "Partner ID"
// @Param partnerInput body domain.PartnerInput true "Partner data"
// @Router /api/v1/partners/update [put]
func (p PartnerController) UpdateOne(ctx *gin.Context) {
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "You are not login!",
		})
		return
	}

	//  Lấy thông tin từ request
	var partnerInput domain.PartnerInput
	if err := ctx.ShouldBindJSON(&partnerInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}
	id := ctx.Query("id")

	err := p.PartnerUseCase.UpdateOne(ctx, id, &partnerInput, fmt.Sprintf("%s", currentUser))
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
