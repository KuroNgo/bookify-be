package activity_log_controller

import (
	"bookify/pkg/shared/constants"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByUserID godoc
// @Summary Get Activity by ID
// @Description Retrieve details of an activity using its user id.
// @Tags Activity Logs
// @Produce json
// @Router /api/v1/activity-logs/get/user_id [get]
func (a ActivityController) GetByUserID(ctx *gin.Context) {
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "You are not login!",
		})
		return
	}

	_, err := a.UserUseCase.GetByID(ctx, fmt.Sprintf("%d", currentUser))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": constants.MsgAPIUnauthorized,
		})
	}

	data, err := a.ActivityUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to get user data: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}
