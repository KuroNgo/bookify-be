package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const MaxFileSize = 5 << 20 // 5MB

func FileUploadMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxFileSize) // Giới hạn request size

		file, header, err := c.Request.FormFile("files")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid file upload",
			})
			return
		}
		defer file.Close()

		// Kiểm tra size của file
		if header.Size > MaxFileSize {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "File size exceeds limit (5MB)",
			})
			return
		}

		c.Set("filePath", header.Filename)
		c.Set("file", file)

		c.Next()
	}
}
