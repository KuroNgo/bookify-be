package middleware

import (
	activitylogcontroller "bookify/internal/api/controller/activity_log"
	"bookify/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"time"
)

// StructuredLogger Ghi log cho các thao tác lỗi
func StructuredLogger(logger *zerolog.Logger, activity *activitylogcontroller.ActivityController) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path

		ctx.Next()

		param := gin.LogFormatterParams{
			TimeStamp: time.Now(),
			Path:      path,
			ClientIP:  ctx.ClientIP(),
			Method:    ctx.Request.Method,
		}

		expireDuration := 30 * 24 * time.Hour
		currentTime := time.Now()
		expireValue := currentTime.Add(expireDuration)

		if ctx.Writer.Status() >= 500 || ctx.Errors != nil || (param.Method == "DELETE" && ctx.Writer.Status() == 200) {
			currentUser, _ := ctx.Get("currentUser")

			param.Latency = time.Since(start).Truncate(time.Millisecond)
			param.StatusCode = ctx.Writer.Status()
			param.ErrorMessage = ctx.Errors.ByType(gin.ErrorTypePrivate).String()

			logger.Error().
				Str("client_id", param.ClientIP).
				Str("method", param.Method).
				Int("status_code", param.StatusCode).
				Int("body_size", ctx.Writer.Size()).
				Str("path", param.Path).
				Str("latency", param.Latency.String()).
				Msg(param.ErrorMessage)

			newLog := &domain.ActivityLogInput{
				ClientIP:     param.ClientIP,
				Method:       param.Method,
				StatusCode:   param.StatusCode,
				BodySize:     ctx.Writer.Size(),
				Path:         path,
				Latency:      param.Latency.String(),
				Error:        param.ErrorMessage,
				ActivityTime: param.TimeStamp,
				ExpireAt:     expireValue,
			}

			go func() {
				err := activity.ActivityUseCase.CreateOne(ctx, newLog, fmt.Sprintf("%s", currentUser))
				if err != nil {
					logger.Error().Err(err).Msg("Failed to create activity log")
					return
				}
			}()

		}
	}
}
