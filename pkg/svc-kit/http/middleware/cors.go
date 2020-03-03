package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS returns cors middleware
func CORS(allowOrigins []string) gin.HandlerFunc {
	c := cors.DefaultConfig()

	c.AllowCredentials = true
	c.AllowOrigins = allowOrigins
	c.AllowHeaders = append(c.AllowHeaders, "Authorization")
	c.AllowMethods = append(c.AllowMethods, "PATCH")
	c.AllowMethods = append(c.AllowMethods, "DELETE")

	return cors.New(c)
}
