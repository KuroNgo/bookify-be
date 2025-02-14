package routes

import (
	"bookify/internal/api/data_seeder"
	"bookify/internal/api/middleware"
	employee_route "bookify/internal/api/routes/employee"
	"bookify/internal/api/routes/event"
	event_ticket_route "bookify/internal/api/routes/event_ticket"
	"bookify/internal/api/routes/event_type"
	organization_route "bookify/internal/api/routes/organization"
	partner_route "bookify/internal/api/routes/partner"
	"bookify/internal/api/routes/user"
	venue_route "bookify/internal/api/routes/venue"
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
	privateRouterV1 := gin.Group("/api/v1")
	userRouter := gin.Group("/api/v1")
	router := gin.Group("")

	publicRouterV1.Use(
		middleware.CORSPublic(),
		middleware.Recover(),
		gzip.Gzip(gzip.DefaultCompression,
			gzip.WithExcludedPaths([]string{",*"})),
	)

	privateRouterV1.Use(
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
	user.UserRouter(env, timeout, db, client, userRouter)
	event.EventsRouter(env, timeout, db, client, publicRouterV1)
	event.AdminEventsRouter(env, timeout, db, client, privateRouterV1)
	event_type.EventTypeRouter(env, timeout, db, publicRouterV1)
	event_type.AdminEventTypeRouter(env, timeout, db, privateRouterV1)
	partner_route.PartnerRouter(env, timeout, db, publicRouterV1)
	partner_route.AdminPartnerRouter(env, timeout, db, privateRouterV1)
	organization_route.OrganizationRouter(env, timeout, db, publicRouterV1)
	organization_route.AdminOrganizationRouter(env, timeout, db, privateRouterV1)
	venue_route.VenueRouter(env, timeout, db, publicRouterV1)
	venue_route.AdminVenueRouter(env, timeout, db, privateRouterV1)
	employee_route.EmployeeRouter(env, timeout, db, privateRouterV1)
	employee_route.AdminEmployeeRouter(env, timeout, db, privateRouterV1)
	event_ticket_route.EventTicketRouter(env, timeout, db, privateRouterV1)
	event_ticket_route.AdminEventTicketRouter(env, timeout, db, privateRouterV1)

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
