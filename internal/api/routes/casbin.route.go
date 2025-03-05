package routes

import (
	"bookify/internal/config"
	casbin_routers "bookify/pkg/interface/casbin/routers"
	"github.com/gin-gonic/gin"
)

func CasbinRouter(env *config.Database, group *gin.RouterGroup) {

	router := group.Group("")
	casbin_routers.CasbinRouter(router, env)
}
