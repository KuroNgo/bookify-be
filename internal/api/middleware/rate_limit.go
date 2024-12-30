package middleware

import (
	"bookify/pkg/shared/constants"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

const (
	maxRequests     = 5
	perMinutePeriod = 15 * time.Second
)

var (
	ipRequestsCounts = make(map[string]int) // can use some distributed db
	mutex            = &sync.Mutex{}
)

func RateLimiter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		mutex.Lock()
		defer mutex.Unlock()
		count := ipRequestsCounts[ip]
		if count >= maxRequests {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":  "fail",
				"message": constants.MsgAPIRateLimitExceeded,
			})
			return
		}

		ipRequestsCounts[ip] = count + 1
		time.AfterFunc(perMinutePeriod, func() {
			mutex.Lock()
			defer mutex.Unlock()

			ipRequestsCounts[ip] = ipRequestsCounts[ip] - 1
		})

		ctx.Next()
	}

}
