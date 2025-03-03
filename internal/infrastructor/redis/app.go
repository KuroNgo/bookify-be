package redis

import (
	"bookify/internal/config"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	Env   *config.Database
	Redis *redis.Client
}

func App() (*Application, *redis.Client) {
	app := &Application{}
	app.Env = config.NewEnv()
	app.Redis = NewRedisClient(app.Env)
	return app, app.Redis
}

func (app *Application) CloseDBConnection() {
	CloseRedisConnection(app.Redis)
}
