package routes

import (
	user_controller "bookify/internal/api/controller/user"
	"bookify/internal/api/middleware"
	"bookify/internal/config"
	"bookify/internal/domain"
	user_repository "bookify/internal/repository/user/repository"
	user_usecase "bookify/internal/usecase/user/usecase"
	"bookify/pkg/interface/cloudinary/middlewares"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func UserRouter(env *config.Database, timeout time.Duration, db *mongo.Database, client *mongo.Client, group *gin.RouterGroup) {
	ur := user_repository.NewUserRepository(db, domain.CollectionUser)

	user := &user_controller.UserController{
		UserUseCase: user_usecase.NewUserUseCase(env, timeout, ur, client),
		Database:    env,
	}

	router := group.Group("/users")
	router.POST("/login", middleware.RateLimiter(), user.LoginUser)
	router.POST("/signup", user.SignUp)
	router.PATCH("/update", middlewares.FileUploadMiddleware(), middleware.DeserializeUser(), user.UpdateUser)
	router.PATCH("/update/image", middleware.DeserializeUser(), user.UpdateImage)
	router.PATCH("/verify", user.VerificationCode)
	router.PATCH("/verify/password", user.VerificationCodeForChangePassword)
	router.PATCH("/password/forget", user.ChangePassword)
	router.POST("/forget", user.ForgetPasswordInUser)
	router.GET("/get/info", user.GetMe)
	router.GET("/get/refresh", user.RefreshToken)
	router.DELETE("/current/delete", middleware.DeserializeUser(), user.DeleteCurrentUser)
	router.GET("/logout", middleware.DeserializeUser(), user.LogoutUser)

	google := group.Group("/auth")
	google.GET("/google/callback", user.GoogleLoginWithUser)
}
