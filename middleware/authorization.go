package middleware

import (
	"github.com/gin-gonic/gin"
)

func (m *Middleware) UserCan(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()
	}
}
