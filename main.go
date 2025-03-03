package main

import (
	"bookify/internal/api/routes"
	"bookify/internal/infrastructor/mongodb"
	"bookify/internal/infrastructor/redis"
	cronjob "bookify/pkg/shared/schedules"
	"github.com/gin-gonic/gin"
	_ "net/http/pprof"
	"time"
)

// @title Bookify
// @version 1.0
// @description This is a server for Kuro API

// @contact.name API Support
// @contact.url
// @contact.email hoaiphong01012002@gmail.com

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Khởi tạo MongoDB
	app, client := mongodb.App()
	env := app.Env
	db := app.MongoDB.Database(env.DBName)
	defer app.CloseDBConnection()

	// Khởi tạo Redis
	appRedis, clientRedis := redis.App()
	envRedis := appRedis.Env
	defer appRedis.CloseDBConnection()

	// Khởi tạo cronjob
	cr := cronjob.NewCronScheduler()

	timeout := time.Duration(env.ContextTimeout) * time.Second
	cacheTTL := time.Minute * 5

	// Khởi tạo Gin
	_gin := gin.Default()

	// Truyền Redis client thay vì `dbRedis`
	routes.SetUp(env, envRedis, cr, timeout, client, clientRedis, db, _gin, cacheTTL)

	err := _gin.Run(env.ServerAddress)
	if err != nil {
		return
	}
}
