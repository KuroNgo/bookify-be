package event_controller

import (
	"bookify/internal/config"
	eventusecase "bookify/internal/usecase/event/usecase"
)

type EventController struct {
	Database     *config.Database
	EventUseCase eventusecase.IEventUseCase
}
