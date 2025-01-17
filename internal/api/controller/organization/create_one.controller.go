package organization_controller

import (
	"bookify/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateOne godoc
// @Summary Create a new organization
// @Description Creates a new organization for the current user
// @Tags Organizations
// @Accept json
// @Produce json
// @Param eventType body domain.OrganizationInput true "Organization Input Body"
// @Router /api/v1/organization/create [post]
func (o OrganizationController) CreateOne(ctx *gin.Context) {
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

	err := o.OrganizationUseCase.CreateOne(ctx, &organizationInput, fmt.Sprintf("%s", currentUser))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	// Trả về phản hồi thành công
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
