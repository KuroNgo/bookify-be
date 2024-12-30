package routes

import (
	"bookify/docs"
	"bookify/internal/api/controller"
	"bookify/internal/api/data_seeder"
	"bookify/internal/api/middleware"
	"bookify/internal/config"
	"bookify/internal/domain"
	"bookify/internal/repository/user/repository"
	"bookify/internal/usecase"
	"bookify/pkg/interface/cloudinary/middlewares"
	"context"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func SetUp(env *config.Database, timeout time.Duration, client *mongo.Client, db *mongo.Database, gin *gin.Engine, cacheTTL time.Duration) {
	publicRouterV1 := gin.Group("/api/v1")
	userRouter := gin.Group("/api/v1")
	router := gin.Group("")

	publicRouterV1.Use(
		middleware.CORSPrivate(),
		middleware.Recover(),
		gzip.Gzip(gzip.DefaultCompression,
			gzip.WithExcludedPaths([]string{",*"})),
		middleware.DeserializeUser(),
	)

	userRouter.Use(
		middleware.CORSPrivate(),
		middleware.Recover(),
		gzip.Gzip(gzip.DefaultCompression,
			gzip.WithExcludedPaths([]string{",*"})),
	)

	// This is a CORS method for check IP validation
	router.OPTIONS("/*path", middleware.OptionMessages)

	SwaggerRouter(env, timeout, db, router)
	UserRouter(env, timeout, db, client, userRouter)

	err := data_seeder.DataSeeds(context.Background(), client)
	if err != nil {
		fmt.Print("data seed is error")
	}

	routeCount := countRoutes(gin)
	fmt.Printf("The number of API endpoints: %d\n", routeCount)
}

func countRoutes(r *gin.Engine) int {
	count := 0
	routes := r.Routes()
	for range routes {
		count++
	}
	return count
}

func init() {
	ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:8080/"),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.DeepLinking(true),
		ginSwagger.PersistAuthorization(true),
	)

	// Save pprof handlers first.
	pprofMux := http.DefaultServeMux
	http.DefaultServeMux = http.NewServeMux()

	// Pprof server.
	go func() {
		fmt.Println(http.ListenAndServe("localhost:8000", pprofMux))
	}()
}

func SwaggerRouter(env *config.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	router := group.Group("")

	docs.SwaggerInfo.BasePath = ""
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//route automatically
	//Thực hiện tự động chuyển hướng khi chạy chương trình
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/swagger/index.html")
	})
}

func UserRouter(env *config.Database, timeout time.Duration, db *mongo.Database, client *mongo.Client, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, domain.CollectionUser)

	user := &controller.UserController{
		UserUseCase: usecase.NewUserUseCase(env, timeout, ur, client),
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
