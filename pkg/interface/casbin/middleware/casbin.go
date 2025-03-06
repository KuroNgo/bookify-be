package middlewares

import (
	user_usecase "bookify/internal/usecase/user/usecase"
	principle "bookify/pkg/interface/casbin/principles"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// rbac.go

func CheckAuth(userUseCase user_usecase.IUserUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUser, exist := ctx.Get("currentUser")
		if !exist {
			fmt.Println("No currentUser in context")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userData, err := userUseCase.GetByID(ctx, fmt.Sprint(currentUser))
		if err != nil {
			fmt.Println("GetByID error:", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		path := "http://localhost:8080" + ctx.Request.URL.Path
		action := ctx.Request.Method // Dùng method thực tế thay vì hardcode "GET"
		fmt.Printf("Checking: Role=%s, Path=%s, Action=%s\n", userData.Role, path, action)

		// In tất cả policy
		policies, err := principle.Rbac.GetPolicy()
		if err != nil {
			fmt.Println("Permission denied")
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}
		fmt.Println("Policies:", policies)
		roles, err := principle.Rbac.GetRolesForUser(userData.Role)
		if err != nil {
			fmt.Println("Permission denied")
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}
		fmt.Println("Roles for", userData.Role, ":", roles)

		ok, err := principle.Rbac.Enforce(userData.Role, path, action)
		if err != nil {
			fmt.Println("Enforce error:", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !ok {
			fmt.Println("Permission denied")
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
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
