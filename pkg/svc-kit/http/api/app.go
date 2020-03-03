package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// App return new app instance.
type App struct {
	Version string
	Engine  *gin.Engine
}

// NewApp return new app instance.
func NewApp(v string) *App {
	return &App{
		Version: v,
		Engine:  gin.Default(),
	}
}

// Run http server.
func (app *App) Run(addr string) error {
	return app.Engine.Run(addr)
}

// GetVersion Get app version
func (app *App) GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"commit": app.Version})
}
