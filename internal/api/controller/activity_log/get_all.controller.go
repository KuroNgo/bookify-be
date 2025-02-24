package activity_log_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAll godoc
// @Summary Get all activity log
// @Description Retrieve a list of all activity log.
// @Tags Activity Logs
// @Produce json
// @Router /api/v1/activity-logs/get/all [get]
func (a ActivityController) GetAll(ctx *gin.Context) {
	data, err := a.ActivityUseCase.GetAll(ctx)
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
