package middlewares

import (
	principle "bookify/pkg/interface/casbin/principles"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// rbac.go

func CheckAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUser, exist := ctx.Get("currentUser")
		if !exist {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		path := "http://localhost:8080" + ctx.Request.URL.Path

		// need check role
		ok, err := principle.Rbac.Enforce(fmt.Sprintf("%s", currentUser), path, "GET")
		if !ok || err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadGateway, "can not get data")
			return
		}

		ctx.Next()
	}
}

func Notification() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusBadGateway, "can not get data")
		return
	}
}
