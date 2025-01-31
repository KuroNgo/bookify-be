package user_controller

import (
	"bookify/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UpdateUserProfileNotImage
// @Summary Update profile user
// @Description Update profile user
// @Tags User
// @Accept json
// @Produce json
// @Param UpdateUserProfileNotImage body domain.UpdateUserSettings true "User data"
// @Security ApiKeyAuth
// @Router /api/v1/users/update/profile/non-image [patch]
func (u *UserController) UpdateUserProfileNotImage(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}

	var input domain.UpdateUserSettings
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	err := u.UserUseCase.UpdateProfileNotImage(ctx, fmt.Sprint(currentUser), &input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Updated user successful",
	})
}
