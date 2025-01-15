package event_controller

import (
	"bookify/internal/config"
	"bookify/internal/usecase/event/usecase"
)

type EventController struct {
	Database     *config.Database
	EventUseCase usecase.IEventUseCase
}
