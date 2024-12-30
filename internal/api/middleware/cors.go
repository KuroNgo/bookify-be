package middleware

import "github.com/gin-gonic/gin"

var (
	host1 = "http://localhost:3000"
	host2 = "http://localhost:5173"
	host3 = "https://bookify-fe.vercel.app"
	host4 = "https://bookify-be-production.up.railway.app"
)

func CORSPublic() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.Request.Header.Get("Origin")
		allowedOrigins := []string{host1, host2, host3, host4}

		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
				ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Content-Length, Accept-Encoding")
				ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE, OPTIONS")

				if ctx.Request.Method == "OPTIONS" {
					ctx.AbortWithStatus(204)
					return
				}
				break
			}
		}

		ctx.Next()
	}
}

func CORSPrivate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.Request.Header.Get("Origin")
		if origin == host2 || origin == host3 || origin == host4 {
			//if origin == host2 {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,Content,Content-Length,Accept-Encoding")
			ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, PUT, PATCH, POST, DELETE, OPTIONS")

			if ctx.Request.Method == "OPTIONS" {
				ctx.AbortWithStatus(204)
				return
			}

			ctx.Next()
		}
	}
}

func OptionMessages(ctx *gin.Context) {
	origin := ctx.Request.Header.Get("Origin")

	if origin == host1 || origin == host2 || origin == host3 || origin == host4 {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,Content-Length")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, PUT, PATCH, POST, DELETE, OPTIONS")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}
