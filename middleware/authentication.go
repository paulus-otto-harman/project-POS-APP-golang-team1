package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/handler"
)

func (m *Middleware) Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		isValid, _ := validateToken(token, m.secretKey)

		if !isValid {
			handler.BadResponse(c, "Unauthorized", http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Next()
	}
}
