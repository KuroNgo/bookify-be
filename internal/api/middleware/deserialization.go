package middleware

import (
	"bookify/internal/infrastructor/mongodb"
	"bookify/pkg/shared/token"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var CacheJWT = make(map[string]interface{}) // Cache tạm thời

func DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string
		cookie, err := ctx.Cookie("access_token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		} else if err == nil {
			accessToken = cookie
		}

		if accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "fail",
				"message": "You are not logged in",
			})
			return
		}

		// 👉 Kiểm tra cache trước khi xác thực token
		if user, found := CacheJWT[accessToken]; found {
			ctx.Set("currentUser", user)
			ctx.Next()
			return
		}

		app, _ := mongodb.App()
		env := app.Env

		sub, err := token.ValidateToken(accessToken, env.AccessTokenPublicKey)
		if err != nil {
			fmt.Println("The err is: ", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "fail",
				"message": err.Error(),
			})
			return
		}

		// Lưu vào cache (nếu có Redis thì lưu vào Redis thay vì map)
		CacheJWT[accessToken] = sub

		ctx.Set("currentUser", sub)
		ctx.Next() // Cho phép tiếp tục các handler khác nếu không có lỗi
	}
}
