package event_type

import (
	"bookify/internal/config"
	event_type_usecase "bookify/internal/usecase/event_type/usecase"
)

type EventController struct {
	Database         *config.Database
	EventTypeUseCase event_type_usecase.IEventTypeRepository
}
