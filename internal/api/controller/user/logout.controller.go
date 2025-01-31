package user_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// LogoutUser
// @Summary Logout user
// @Description Logout the current user
// @Tags User
// @Accept  json
// @Produce  json
// @Router /api/v1/users/logout [get]
func (u *UserController) LogoutUser(ctx *gin.Context) {
	//_, exists := ctx.Get("currentUser")
	//if !exists {
	//	ctx.JSON(http.StatusUnauthorized, gin.H{
	//		"status":  "fail",
	//		"message": "You are not logged in!",
	//	})
	//	return
	//}
	//
	//user, err := u.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	//if err != nil || user.ID == primitive.NilObjectID {
	//	ctx.JSON(http.StatusUnauthorized, gin.H{
	//		"status":  "Unauthorized",
	//		"message": "You are not authorized to perform this action!",
	//	})
	//	return
	//}

	ctx.SetCookie("access_token", "", -1, "/", u.Database.ClientServer, false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", u.Database.ClientServer, false, true)
	ctx.SetCookie("is_logged", "", -1, "/", u.Database.ClientServer, false, false)

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
