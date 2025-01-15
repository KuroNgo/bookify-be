package user_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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
