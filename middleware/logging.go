package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		log.Printf("Started %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
		log.Printf("Completed %s %s in %v", c.Request.Method, c.Request.URL.Path, time.Since(start))
	}
}
