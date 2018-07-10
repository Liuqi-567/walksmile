package controllers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
)

func Domain() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(c.Request.Host, "localhost") && !strings.Contains(c.Request.Host, "walksmile.com") {
			c.AbortWithStatus(404)
		}
		c.Next()
	}
}

func IsMobile() gin.HandlerFunc {
	return func(c *gin.Context) {
		ua := user_agent.New(c.Request.UserAgent())
		if ua.Mobile() {
			c.Set("isMobile", true)
		}
		c.Next()
	}
}
