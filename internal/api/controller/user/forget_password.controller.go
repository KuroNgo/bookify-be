package user_controller

import (
	"bookify/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ForgetPasswordInUser allows the user to request a password reset.
// @Summary User forget password
// @Description Sends an email with a verification code for password reset
// @Tags User
// @Accept  json
// @Produce  json
// @Param forgetInput body domain.ForgetPassword true "Forget password input"
// @Router /api/v1/users/forget [post]
func (u *UserController) ForgetPasswordInUser(ctx *gin.Context) {
	var forgetInput domain.ForgetPassword
	if err := ctx.ShouldBindJSON(&forgetInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	err := u.UserUseCase.ForgetPassword(ctx, forgetInput.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	// Trả về phản hồi thành công
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "We sent an email with a verification code to your email",
	})
}

// VerificationCodeForChangePassword verifies the code sent to the user for changing password.
// @Summary Verify code for password change
// @Description Verifies the code sent to the user's email for changing password
// @Tags User
// @Accept  json
// @Produce  json
// @Param verificationCode body domain.VerificationInput true "Verification code input"
// @Router /api/v1/users/verify/password [patch]
func (u *UserController) VerificationCodeForChangePassword(ctx *gin.Context) {
	var verificationCode domain.VerificationInput
	if err := ctx.ShouldBindJSON(&verificationCode); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	err := u.UserUseCase.UpdateVerifyForChangePassword(ctx, verificationCode.VerificationCode)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Set cookie
	ctx.SetCookie("verification_code", verificationCode.VerificationCode, 0, "/", u.Database.ClientServer, false, true)

	// Trả về phản hồi thành công
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// ChangePassword allows the user to change their password using a verification code.
// @Summary Change password
// @Description Allows the user to change their password after verifying the code
// @Tags User
// @Accept  json
// @Produce  json
// @Param changePasswordInput body domain.ChangePasswordInput true "Change password input"
// @Router /api/v1/users/password/forget [patch]
func (u *UserController) ChangePassword(ctx *gin.Context) {
	cookie, err := ctx.Cookie("verification_code")
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "Verification code is missing!",
		})
		return
	}

	var input domain.ChangePasswordInput
	if err = ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	err = u.UserUseCase.UpdatePassword(ctx, cookie, &input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	ctx.SetCookie("verification_code", "", -1, "/", u.Database.ClientServer, false, true)

	// Trả về phản hồi thành công
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
