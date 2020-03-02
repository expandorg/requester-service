package endpoints

import (
	"github.com/expandorg/requester-service/pkg/drafttemplates"
	"github.com/expandorg/requester-service/pkg/http/app/ctx"
	"github.com/gemsorg/svc-kit/http/api"
	"github.com/gin-gonic/gin"
)

// TemplatesRoutes setup templates routes
func TemplatesRoutes(r *gin.RouterGroup, h *api.Handler, templates drafttemplates.Service) {
	c := templatesController{
		templates: templates,
	}

	r.GET("/tasks/templates", h.Handle(c.GetTaskTemplates))
	r.GET("/tasks/templates/:templateID", h.Handle(c.GetTaskTemplate))
}

type templatesController struct {
	templates drafttemplates.Service
}

func (tc *templatesController) GetTaskTemplates(c *gin.Context) (interface{}, error) {
	templates, err := tc.templates.List(ctx.ReqCtx(c))
	if err != nil {
		return nil, err
	}
	return gin.H{"templates": templates}, nil
}

func (tc *templatesController) GetTaskTemplate(c *gin.Context) (interface{}, error) {
	template, err := tc.templates.FindByID(ctx.ReqCtx(c), ctx.TemplateID(c))
	if err != nil {
		return nil, err
	}
	return gin.H{"template": template}, nil
}
