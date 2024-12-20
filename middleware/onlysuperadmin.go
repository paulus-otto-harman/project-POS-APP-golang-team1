package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"project/handler"
)

func (m *Middleware) OnlySuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {

		role, err := m.cacher.HGet(fmt.Sprintf("user:%s", c.GetString("user-id")), "role")
		if err != nil {
			handler.BadResponse(c, "Unauthorized", http.StatusUnauthorized)
			c.Abort()
			return
		}

		if role != "super admin" {
			handler.BadResponse(c, "Forbidden", http.StatusForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}
