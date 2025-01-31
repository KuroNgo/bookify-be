package user_controller

import (
	"bookify/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// UpdateUserProfile godoc
// @Summary Update user profile
// @Description Update the profile information of the currently logged-in user
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param full_name formData string false "Full name of the user"
// @Param gender formData string false "Gender of the user"
// @Param vocation formData string false "Vocation or profession"
// @Param address formData string false "Address of the user"
// @Param city formData string false "City of residence"
// @Param region formData string false "Region or state"
// @Param date_of_birth formData string false "Date of birth in YYYY-MM-DD format"
// @Param show_interest formData boolean false "Show interests in profile (true/false)"
// @Param social_media formData boolean false "Enable social media sharing (true/false)"
// @Param file formData file false "Profile picture file"
// @Router /api/v1/user/update/profile [patch]
func (u *UserController) UpdateUserProfile(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}

	fullName := ctx.Request.FormValue("full_name")
	phone := ctx.Request.FormValue("phone")
	gender := ctx.Request.FormValue("gender")
	vocation := ctx.Request.FormValue("vocation")
	address := ctx.Request.FormValue("address")
	city := ctx.Request.FormValue("city")
	region := ctx.Request.FormValue("region")
	dateOfBirth := ctx.Request.FormValue("date_of_birth")
	showInterest := ctx.Request.FormValue("show_interest")
	socialMedia := ctx.Request.FormValue("social_media")

	parseShowInterest, _ := strconv.ParseBool(showInterest)
	parseSocialMedia, _ := strconv.ParseBool(socialMedia)

	input := domain.UpdateUserSettings{
		Gender:       gender,
		Phone:        phone,
		Vocation:     vocation,
		Address:      address,
		City:         city,
		Region:       region,
		DateOfBirth:  dateOfBirth,
		FullName:     fullName,
		ShowInterest: parseShowInterest,
		SocialMedia:  parseSocialMedia,
	}

	file, _ := ctx.FormFile("file")
	// Delete image in cloudinary before save by assetID to update user info
	imageURL, err := u.UserUseCase.UpdateProfile(ctx, fmt.Sprint(currentUser), &input, file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"image_url": imageURL,
		"message":   "Updated user successful",
	})
}
