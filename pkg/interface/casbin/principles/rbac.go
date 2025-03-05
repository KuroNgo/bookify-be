package principle

import (
	"bookify/internal/config"
	"fmt"
	"github.com/casbin/casbin/v2"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"log"
)

var Rbac *casbin.Enforcer

func SetUp(env *config.Database) *casbin.Enforcer {
	var mongodbURI string
	if env.DBUser != "" && env.DBPassword != "" {
		mongodbURI = fmt.Sprintf("mongodb+srv://%s:%s@andrew.8ulkv.mongodb.net/?retryWrites=true&w=majority", env.DBUser, env.DBPassword)
	} else {
		mongodbURI = fmt.Sprintf("mongodb://%s:%s", env.DBHost, env.DBPort)
	}
	a, err := mongodbadapter.NewAdapter(mongodbURI)
	if err != nil {
		log.Fatalln(err)
	}

	r, err := casbin.NewEnforcer("./pkg/interface/casbin/config/rbac_model.conf", a)
	if err != nil {
		log.Fatalln(err)
	}

	err = r.LoadPolicy()
	if err != nil {
		return nil
	}

	Rbac = r

	return r
}
