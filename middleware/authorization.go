package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"project/helper"
)

func (m *Middleware) CanAccess(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(c.GetString("user-id"))
		userID := c.GetString("user-id")
		role, err := m.cacher.HGet("user:"+userID, "role")
		if err != nil {
			c.Abort()
			return
		}
		log.Println("role", role)

		if role == "super admin" {
			c.Next()
			return
		}

		log.Println("requires permission", permission)
		var isAuthorized bool
		isAuthorized, err = m.cacher.SIsMember(fmt.Sprintf("user:%s:permission", userID), permission)

		if err != nil {
			helper.BadResponse(c, "server error", http.StatusInternalServerError)
			c.Abort()
			return
		}

		if !isAuthorized {
			helper.BadResponse(c, "unauthorized", http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Next()
	}
}
