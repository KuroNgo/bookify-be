package infrastructor

import (
	"bookify/internal/config"
	mongodriven "go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	Env     *config.Database
	MongoDB *mongodriven.Client
}

func App() (*Application, *mongodriven.Client) {
	app := &Application{}
	app.Env = config.NewEnv()
	app.MongoDB = NewMongoDatabase(app.Env)
	return app, app.MongoDB
}

func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.MongoDB)
}
