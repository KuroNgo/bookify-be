package activity_log_controller

import (
	"bookify/internal/config"
	activity_log_usecase "bookify/internal/usecase/activity_log/usecase"
	user_usecase "bookify/internal/usecase/user/usecase"
)

type ActivityController struct {
	ActivityUseCase activity_log_usecase.IActivityUseCase
	UserUseCase     user_usecase.IUserUseCase
	Database        *config.Database
}
