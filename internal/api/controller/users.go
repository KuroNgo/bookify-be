package controller

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	"bookify/internal/usecase"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	Database    *config.Database
	UserUseCase usecase.IUserUseCase
}

// DeleteCurrentUser delete the user's information
// @Summary Delete User Information
// @Description Deletes the user's information
// @Tags User
// @Accept json
// @Produce json
// @Router /api/v1/users/delete [delete]
// @Security CookieAuth
func (u *UserController) DeleteCurrentUser(c *gin.Context) {
	currentUser, exist := c.Get("currentUser")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "You are not login!",
		})
		return
	}

	err := u.UserUseCase.DeleteOne(c, fmt.Sprint(currentUser))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

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

// GetMe retrieves the user information based on the access token.
// @Summary Get User Information
// @Description Retrieves the user's information using the access token stored in cookies.
// @Tags User
// @Accept  json
// @Produce  json
// @Router /api/v1/users/get/info [get]
// @Security CookieAuth
func (u *UserController) GetMe(ctx *gin.Context) {
	// Lấy cookie access_token từ request
	cookie, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}

	// Gọi use case để xử lý logic nghiệp vụ
	result, err := u.UserUseCase.GetByIDForCheckCookie(ctx, cookie)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "Failed to get user data: " + err.Error(),
		})
		return
	}

	// Trả về phản hồi thành công
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

// LoginUser
// @Summary Login user
// @Description Login user
// @Tags User
// @Accept json
// @Produce json
// @Param LoginUser body domain.SignIn true "User data"
// @Security ApiKeyAuth
// @Router /api/v1/users/login [post]
func (l *UserController) LoginUser(ctx *gin.Context) {
	//  Lấy thông tin từ request
	var userInput domain.SignIn
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	// Xử lý logic nghiệp vụ và tìm kiếm người dùng
	user, err := l.UserUseCase.LoginUser(ctx, &userInput)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.SetCookie("access_token", user.AccessToken, 0, "/", l.Database.ClientServer, false, true)
	ctx.SetCookie("refresh_token", user.RefreshToken, 0, "/", l.Database.ClientServer, false, true)
	ctx.SetCookie("is_logged", user.IsLogged, 0, "/", l.Database.ClientServer, false, false)

	// Trả về phản hồi thành công
	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"data":      user,
		"is_logged": user.IsLogged,
	})
}

// GoogleLoginWithUser
// @Summary Login Google
// @Description  Login the user's google, but the function not use with swagger.
// @Tags User
// @Router /api/v1/users/google/callback [get]
func (u *UserController) GoogleLoginWithUser(c *gin.Context) {
	code := c.Query("code")

	userData, response, err := u.UserUseCase.LoginGoogle(c, code)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	c.SetCookie("access_token", response.AccessToken, 0, "/", "localhost", false, true)
	c.SetCookie("refresh_token", response.RefreshToken, 0, "/", "localhost", false, true)
	c.SetCookie("is_logged", response.IsLogged, 0, "/", "localhost", false, false)

	c.JSON(http.StatusOK, gin.H{
		"token": response.SignedToken,
		"user":  userData.User,
	})
}

// RefreshToken refreshes the user's access token.
// @Summary Refresh Access Token
// @Description Refresh the user's access token using a valid refresh token stored in cookies.
// @Tags User
// @Accept  json
// @Produce  json
// @Router /api/v1/users/get/refresh [get]
// @Security CookieAuth
func (u *UserController) RefreshToken(ctx *gin.Context) {
	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "could not refresh access token",
		})
		return
	}

	data, err := u.UserUseCase.RefreshToken(ctx, cookie)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "could not refresh access token",
		})
		return
	}

	ctx.SetCookie("access_token", data.AccessToken, 0, "/", u.Database.ClientServer, false, true)
	ctx.SetCookie("refresh_token", data.RefreshToken, 0, "/", u.Database.ClientServer, false, true)
	ctx.SetCookie("is_logged", data.IsLogged, 0, "/", u.Database.ClientServer, false, false)

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"is_logged": data.IsLogged,
	})
}

// LogoutUser
// @Summary Logout user
// @Description Logout the current user
// @Tags User
// @Accept  json
// @Produce  json
// @Router /api/v1/users/logout [get]
func (u *UserController) LogoutUser(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}

	user, err := u.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	ctx.SetCookie("access_token", "", -1, "/", u.Database.ClientServer, false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", u.Database.ClientServer, false, true)
	ctx.SetCookie("is_logged", "", -1, "/", u.Database.ClientServer, false, false)

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// SignUp Create a new user
// @Summary Register user
// @Description Register a new user with form data
// @Tags User
// @Accept x-www-form-urlencoded
// @Accept multipart/form-data
// @Produce json
// @Param email formData string true "Email of the user" example("john.doe@example.com")
// @Param password formData string true "Password of the user" example("securepassword123")
// @Param fullName formData string true "Full name of the user" example("John Doe")
// @Param avatarUrl formData string false "Avatar URL of the user" example("http://example.com/avatar.jpg")
// @Param phone formData string true "Phone number of the user" example("+1234567890")
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

// UpdateImage updates the user's information
// @Summary Update User Information
// @Description Updates the user's first name, last name, and username
// @Tags User
// @Accept json
// @Produce json
// @Param file formData file false "Image file to upload"
// @Router /api/v1/users/update/image [put]
// @Security CookieAuth
func (u *UserController) UpdateImage(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}

	file, _ := ctx.FormFile("file")
	err := u.UserUseCase.UpdateImage(ctx, fmt.Sprint(currentUser), file)
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
