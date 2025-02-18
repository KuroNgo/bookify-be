package user_controller

import (
	"bookify/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SignUp Create a new user
// @Summary Register user
// @Description Register a new user with form data
// @Tags User
// @Produce json
// @Param SignUp body domain.SignupUser true "User data"
// @Security ApiKeyAuth
// @Router /api/v1/users/signup [post]
func (u *UserController) SignUp(ctx *gin.Context) {
	//  Lấy thông tin từ request
	var userInput domain.SignupUser
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	err := u.UserUseCase.SignUp(ctx, &userInput)
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

// VerificationCode Create a new user
// @Summary Register user
// @Description Register a new user with form data
// @Tags User
// @Accept json
// @Produce json
// @Param User body domain.VerificationInput true "User data"
// @Security ApiKeyAuth
// @Router /api/v1/users/verify [patch]
func (u *UserController) VerificationCode(ctx *gin.Context) {
	var verificationCode domain.VerificationInput
	if err := ctx.ShouldBindJSON(&verificationCode); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	user, err := u.UserUseCase.GetByVerificationCode(ctx, verificationCode.VerificationCode)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	// Trả về phản hồi thành công
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user,
	})
}
