package controller

import (
	"bookify/internal/config"
	"bookify/internal/domain"
	"bookify/internal/usecase"
	"bookify/pkg/shared/constants"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type EventController struct {
	Database     *config.Database
	EventUseCase usecase.IEventUseCase
}

// GetByID
// @Summary Get event by ID
// @Description Get details of an event by its ID
// @Tags Events
// @Accept json
// @Produce json
// @Param id query string true "Event ID"
// @Router /api/v1/events/get/id [get]
func (e *EventController) GetByID(ctx *gin.Context) {
	id := ctx.Query("id")

	data, err := e.EventUseCase.GetByID(ctx, id)
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

// GetByStartTime
// @Summary Get events by start time
// @Description Get events starting from a specific date
// @Tags Events
// @Accept json
// @Produce json
// @Param startDate query string true "Start date in YYYY-MM-DD format"
// @Router /api/v1/events/get/start-time [get]
func (e *EventController) GetByStartTime(ctx *gin.Context) {
	startDate := ctx.Query("startDate")

	data, err := e.EventUseCase.GetByStartTime(ctx, startDate)
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

// GetByStartTimePagination
// @Summary Get events by start time with pagination
// @Description Get paginated events starting from a specific date
// @Tags Events
// @Accept json
// @Produce json
// @Param startDate query string true "Start date in YYYY-MM-DD format"
// @Param page query string false "Page number, default is 1"
// @Router /api/v1/events/get/start-time/pagination [get]
func (e *EventController) GetByStartTimePagination(ctx *gin.Context) {
	startDate := ctx.Query("startDate")
	page := ctx.DefaultQuery("page", "1")

	data, pageOutput, currentPage, err := e.EventUseCase.GetByStartTimePagination(ctx, startDate, page)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to get user data: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"data":         data,
		"page":         pageOutput,
		"current_page": currentPage,
	})
}

// GetAll
// @Summary Get all events
// @Description Get a list of all events
// @Tags Events
// @Accept json
// @Produce json
// @Router /api/v1/events/get/all [get]
func (e *EventController) GetAll(ctx *gin.Context) {
	data, err := e.EventUseCase.GetAll(ctx)
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

// GetAllPagination
// @Summary Get all events with pagination
// @Description Get a paginated list of all events
// @Tags Events
// @Accept json
// @Produce json
// @Param page query string false "Page number, default is 1"
// @Router /api/v1/events/get-all/pagination [get]
func (e *EventController) GetAllPagination(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")

	data, pageOutput, currentPage, err := e.EventUseCase.GetAllPagination(ctx, page)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to get user data: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"data":         data,
		"page":         pageOutput,
		"current_page": currentPage,
	})
}

// CreateOne
// @Summary Create a new event
// @Description Add a new event to the system
// @Tags Events
// @Accept json
// @Produce json
// @Param event body domain.EventInput true "Event input data"
// @Router /api/v1/events/create [post]
func (e *EventController) CreateOne(ctx *gin.Context) {
	// Retrieve the current user from the context
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": constants.MsgAPIUnauthorized,
		})
		return
	}

	// Bind input data
	var event domain.EventInput
	event.UserID = fmt.Sprintf("%s", currentUser)
	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data: " + err.Error(),
		})
		return
	}

	// Call the use case to create the event
	if err := e.EventUseCase.CreateOne(ctx, &event); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to create event: " + err.Error(),
		})
		return
	}

	// Respond with success
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Event created successfully",
	})
}

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
	currentUser, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": constants.MsgAPIUnauthorized,
		})
		return
	}

	// Bind input data
	var event domain.EventInput
	event.UserID = fmt.Sprintf("%s", currentUser)
	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data: " + err.Error(),
		})
		return
	}
	id := ctx.Query("id")

	err := e.EventUseCase.UpdateOne(ctx, id, &event)
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

// DeleteOne
// @Summary Delete an event
// @Description Delete an event by its ID
// @Tags Events
// @Accept json
// @Produce json
// @Param id query string true "Event ID"
// @Router /api/v1/events/delete-one [delete]
func (e *EventController) DeleteOne(ctx *gin.Context) {
	_, exist := ctx.Get("currentUser")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": constants.MsgAPIUnauthorized,
		})
		return
	}

	id := ctx.Query("id")

	err := e.EventUseCase.DeleteOne(ctx, id)
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
