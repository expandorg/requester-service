package endpoints

import (
	"github.com/expandorg/requester-service/pkg/dashboard"
	"github.com/expandorg/requester-service/pkg/http/app/ctx"
	"github.com/expandorg/requester-service/pkg/svc-kit/http/api"

	"github.com/gin-gonic/gin"
)

// TasksRoutes setup tasks routes
func TasksRoutes(r *gin.RouterGroup, h *api.Handler, tasks dashboard.Service) {
	c := tasksController{
		tasks: tasks,
	}

	r.GET("/tasks/status/", h.Handle(c.GetTasks))
	r.GET("/tasks/status/:taskStatus", h.Handle(c.GetTasksByStatus))
}

type tasksController struct {
	tasks dashboard.Service
}

func (tc *tasksController) GetTasks(c *gin.Context) (interface{}, error) {
	tasks, err := tc.tasks.ListByUserID(ctx.ReqCtx(c), ctx.UserID(c))
	if err != nil {
		return nil, err
	}
	return gin.H{"tasks": tasks}, nil
}

func (tc *tasksController) GetTasksByStatus(c *gin.Context) (interface{}, error) {
	status, err := ctx.TaskStatus(c)
	if err != nil {
		return nil, api.BadRequest(err)
	}
	tasks, err := tc.tasks.ListByUserIDAndStatus(ctx.ReqCtx(c), ctx.UserID(c), status)
	if err != nil {
		return nil, err
	}
	return gin.H{"tasks": tasks}, nil
}
