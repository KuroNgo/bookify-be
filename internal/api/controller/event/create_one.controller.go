package event_controller

import (
	"github.com/gin-gonic/gin"
)

// CreateOne
// @Summary Create a new event
// @Description Add a new event to the system
// @Tags Events
// @Accept json
// @Produce json
// @Param event body domain.EventInput true "Event input data"
// @Router /api/v1/events/create [post]
func (e *EventController) CreateOne(ctx *gin.Context) {

}
