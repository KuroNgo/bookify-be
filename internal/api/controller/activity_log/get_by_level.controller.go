package activity_log_controller

import (
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
	id := ctx.Query("id")
	data, err := a.ActivityUseCase.GetByID(ctx, id)
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
