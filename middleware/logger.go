package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func (m *Middleware) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		log.Printf("Request %s - %s - %s", c.Request.Method, c.Request.URL.Path, latency)
	}
}
