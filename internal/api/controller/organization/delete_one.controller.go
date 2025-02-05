package organization_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// DeleteOne godoc
// @Summary Delete an organization
// @Description Deletes a specific organization by ID. User must be logged in.
// @Tags Organizations
// @Produce json
// @Param id query string true "Organization ID"
// @Router /api/v1/organizations/delete [delete]
// @Security ApiKeyAuth
func (o OrganizationController) DeleteOne(ctx *gin.Context) {
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "You are not login!",
		})
		return
	}

	id := ctx.Query("id")
	err := o.OrganizationUseCase.DeleteOne(ctx, id, fmt.Sprintf("%s", currentUser))
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
