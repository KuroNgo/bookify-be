package user_controller

import (
	"bookify/internal/config"
	"bookify/internal/usecase/user/usecase"
)

type UserController struct {
	Database    *config.Database
	UserUseCase user_usecase.IUserUseCase
}
