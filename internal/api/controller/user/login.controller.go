package user_controller

import (
	"bookify/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
	//ctx.SetCookie("is_logged", user.IsLogged, 0, "/", l.Database.ClientServer, false, false)

	// Trả về phản hồi thành công
	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"is_logged": user.IsLogged,
		"data":      user,
	})
}
