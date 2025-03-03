package activity_log_controller

import (
	"bookify/internal/config"
	activitylogusecase "bookify/internal/usecase/activity_log/usecase"
	userusecase "bookify/internal/usecase/user/usecase"
)

type ActivityController struct {
	ActivityUseCase activitylogusecase.IActivityUseCase
	UserUseCase     userusecase.IUserUseCase
	Database        *config.Database
}
