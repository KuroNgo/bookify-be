package routes

import (
	"bookify/docs"
	"bookify/internal/config"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

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
