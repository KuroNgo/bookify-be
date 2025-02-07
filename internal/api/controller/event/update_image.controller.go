package event_controller

import (
	"bookify/pkg/shared/constants"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UpdateImage
// @Summary Update an event
// @Description Update details of an existing event
// @Tags Events
// @Accept json
// @Produce json
// @Param id query string true "Event ID"
// @Param file formData file false "Image file to upload"
// @Router /api/v1/event/update/image [patch]
func (e *EventController) UpdateImage(ctx *gin.Context) {
	_, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": constants.MsgAPIUnauthorized,
		})
		return
	}

	//  Lấy thông tin từ request
	id := ctx.Query("id")
	file, _ := ctx.FormFile("file")

	err := e.EventUseCase.UpdateImage(ctx, id, file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to get user data: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
