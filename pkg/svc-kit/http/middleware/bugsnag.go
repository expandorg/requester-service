package middleware

import (
	"os"

	bugsnag "github.com/bugsnag/bugsnag-go"
	bugsnaggin "github.com/bugsnag/bugsnag-go/gin"
	"github.com/gin-gonic/gin"
)

// Bugsnag middleware
func Bugsnag() gin.HandlerFunc {
	return bugsnaggin.AutoNotify(bugsnag.Configuration{
		APIKey:       os.Getenv("BUGSNAG_KEY"),
		ReleaseStage: "production",
	})
}
