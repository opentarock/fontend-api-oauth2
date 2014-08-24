package middleware

import "github.com/gin-gonic/gin"

func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "no-store")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Next()
	}
}
