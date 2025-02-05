package organization_controller

import (
	"bookify/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UpdateOne godoc
// @Summary Update an organization
// @Description Update details of an organization by its ID.
// @Tags Organizations
// @Accept json
// @Produce json
// @Param id query string true "Organization ID"
// @Param body domain.OrganizationInput true "Organization data"
// @Router /api/v1/organizations/update [put]
func (o OrganizationController) UpdateOne(ctx *gin.Context) {
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "You are not login!",
		})
		return
	}

	//  Lấy thông tin từ request
	var organizationInput domain.OrganizationInput
	if err := ctx.ShouldBindJSON(&organizationInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}
	id := ctx.Query("id")

	err := o.OrganizationUseCase.UpdateOne(ctx, id, &organizationInput, fmt.Sprintf("%s", currentUser))
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
