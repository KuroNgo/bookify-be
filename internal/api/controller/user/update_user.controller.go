package user_controller

import (
	"bookify/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UpdateUser updates the user's information
// @Summary Update User Information
// @Description Updates the user's first name, last name, and username
// @Tags User
// @Accept json
// @Produce json
// @Router /api/v1/users/update [put]
// @Security CookieAuth
func (u *UserController) UpdateUser(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}

	fullName := ctx.Request.FormValue("full_name")
	gender := ctx.Request.FormValue("gender")
	vocation := ctx.Request.FormValue("vocation")
	address := ctx.Request.FormValue("address")
	city := ctx.Request.FormValue("city")
	region := ctx.Request.FormValue("region")
	dateOfBirth := ctx.Request.FormValue("data_of_birth")
	input := domain.UpdateUserInfo{
		FullName:    fullName,
		Gender:      gender,
		Vocation:    vocation,
		Address:     address,
		City:        city,
		Region:      region,
		DateOfBirth: dateOfBirth,
	}

	file, _ := ctx.FormFile("file")
	err := u.UserUseCase.UpdateUserInfoOne(ctx, fmt.Sprint(currentUser), &input, file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Updated user",
	})

}
