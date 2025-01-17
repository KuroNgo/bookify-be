package organization_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByID godoc
// @Summary Get organization by ID
// @Description Retrieve details of an organization using its ID.
// @Tags Organizations
// @Produce json
// @Param id query string true "Organization ID"
// @Router /organizations/get/id [get]
func (o OrganizationController) GetByID(ctx *gin.Context) {
	id := ctx.Query("id")
	data, err := o.OrganizationUseCase.GetByID(ctx, id)
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
