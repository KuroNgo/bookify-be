package organization_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAll godoc
// @Summary Get all organizations
// @Description Retrieve a list of all organizations.
// @Tags Organizations
// @Produce json
// @Router /organizations/get/all [get]
func (o OrganizationController) GetAll(ctx *gin.Context) {
	data, err := o.OrganizationUseCase.GetAll(ctx)
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
