package activity_log_controller

import (
	"bookify/pkg/shared/constants"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByLevel godoc
// @Summary Get organization by level
// @Description Retrieve details of an organization using its Level.
// @Tags Activity Logs
// @Produce json
// @Param level query string true "Activity Level"
// @Router /api/v1/activity-logs/get/level [get]
func (a ActivityController) GetByLevel(ctx *gin.Context) {
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

	level := ctx.Query("level")
	data, err := a.ActivityUseCase.GetByLevel(ctx, level)
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
