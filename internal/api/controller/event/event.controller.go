package event_controller

import (
	"bookify/internal/config"
	event_usecase "bookify/internal/usecase/event/usecase"
)

type EventController struct {
	Database     *config.Database
	EventUseCase event_usecase.IEventUseCase
}
