package activity_log_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetByID godoc
// @Summary Get Activity by ID
// @Description Retrieve details of an activity using its id.
// @Tags Activity Logs
// @Produce json
// @Param id query string true "Activity Log ID"
// @Router /api/v1/activity-logs/get/id [get]
func (a ActivityController) GetByID(ctx *gin.Context) {
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
