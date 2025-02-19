package main

import (
	"bookify/internal/api/routes"
	"bookify/internal/infrastructor"
	cronjob "bookify/pkg/shared/cron"
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

	app, client := infrastructor.App()
	env := app.Env
	db := app.MongoDB.Database(env.DBName)
	defer app.CloseDBConnection()

	cr := cronjob.NewCronScheduler()

	timeout := time.Duration(env.ContextTimeout) * time.Second
	cacheTTL := time.Minute * 5

	_gin := gin.Default()

	routes.SetUp(env, cr, timeout, client, db, _gin, cacheTTL)
	err := _gin.Run(env.ServerAddress)
	if err != nil {
		return
	}
}
