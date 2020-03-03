package endpoints

import (
	"github.com/expandorg/requester-service/pkg/http/app/ctx"
	"github.com/expandorg/requester-service/pkg/onboardingtemplate"
	"github.com/expandorg/requester-service/pkg/svc-kit/http/api"
	"github.com/gin-gonic/gin"
)

// OnboardingRoutes setup templates routes
func OnboardingRoutes(r *gin.RouterGroup, h *api.Handler, templates onboardingtemplate.Service) {
	c := onboardingController{
		templates: templates,
	}
	r.GET("/forms/templates", h.Handle(c.GetFormsTemplates))
	r.GET("/forms/templates/:templateID", h.Handle(c.GetFormsTemplate))
}

type onboardingController struct {
	templates onboardingtemplate.Service
}

func (o *onboardingController) GetFormsTemplates(c *gin.Context) (interface{}, error) {
	templates, err := o.templates.List(ctx.ReqCtx(c))
	if err != nil {
		return nil, err
	}
	return gin.H{"templates": templates}, nil
}

func (o *onboardingController) GetFormsTemplate(c *gin.Context) (interface{}, error) {
	template, err := o.templates.FindByID(ctx.ReqCtx(c), ctx.TemplateID(c))
	if err != nil {
		return nil, err
	}
	return gin.H{"template": template}, nil
}
