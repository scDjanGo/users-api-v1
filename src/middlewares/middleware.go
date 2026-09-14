package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LimitBody(n int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, n)
		c.Next()
	}
}

func MediaHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")

		c.Header("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; sandbox")

		c.Header("Cache-Control", "public, max-age=31536000, immutable")

		c.Next()
	}
}
