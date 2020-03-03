package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler wrap all api request
type Handler struct {
	LogError func(c *gin.Context, e *Error)
}

// NewHandler returns new instance of ApiHandler
func NewHandler(logError logFn) *Handler {
	return &Handler{LogError: logError}
}

// HandleError error
func (h *Handler) HandleError(c *gin.Context, e *Error) {
	h.LogError(c, e)
	c.AbortWithStatusJSON(e.Status, gin.H{"error": e.GetPublicMessage()})
}

// Handle Controller handler
func (h *Handler) Handle(method apiCall) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := method(c)
		if err != nil {
			h.HandleError(c, WrapErr(err))
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

// HandleEmpty Controller handler
func (h *Handler) HandleEmpty(method apiCall) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := method(c)
		if err != nil {
			h.HandleError(c, WrapErr(err))
			return
		}
	}
}

// Ensure middleware handler
func (h *Handler) Ensure(fn middlewareCall) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := fn(c)
		if err != nil {
			h.HandleError(c, WrapErr(err))
		}
	}
}

// ControllerFunc Controller method
type apiCall func(*gin.Context) (interface{}, error)

// MiddlewareFunc returns Middleware method
type middlewareCall func(*gin.Context) error

// Logger
type logFn func(c *gin.Context, e *Error)
