package event_controller

import (
	"github.com/gin-gonic/gin"
)

// UpdateOne
// @Summary Update an event
// @Description Update details of an existing event
// @Tags Events
// @Accept json
// @Produce json
// @Param id query string true "Event ID"
// @Param event body domain.EventInput true "Event input data"
// @Router /api/v1/events/update [put]
func (e *EventController) UpdateOne(ctx *gin.Context) {

}
