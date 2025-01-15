package routes

import (
	"bookify/internal/api/data_seeder"
	"bookify/internal/api/middleware"
	"bookify/internal/config"
	"context"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
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
	EventsRouter(env, timeout, db, client, publicRouterV1)

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
