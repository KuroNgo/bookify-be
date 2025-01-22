package user_controller

import (
	"bookify/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UpdateSocialUser docs
// @Summary Update social media links of the current user
// @Description Allows the logged-in user to update their social media links
// @Tags User
// @Accept json
// @Produce json
// @Param body body domain.UpdateSocialMedia true "Social media update payload"
// @Router /api/v1/users/update/social [patch]
func (u *UserController) UpdateSocialUser(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}

	//  Lấy thông tin từ request
	var userInput domain.UpdateSocialMedia
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	err := u.UserUseCase.UpdateSocialMedia(ctx, fmt.Sprint(currentUser), &userInput)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Updated Social user",
	})
}
