package organization_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByUserID godoc
// @Summary Get organization by ID
// @Description Retrieve details of an organization using its ID.
// @Tags Organizations
// @Produce json
// @Router /api/v1/organizations/get/user_id [get]
func (o OrganizationController) GetByUserID(ctx *gin.Context) {
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "You are not login!",
		})
		return
	}

	data, err := o.OrganizationUseCase.GetByUserID(ctx, fmt.Sprintf("%s", currentUser))
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
