package drafttemplates

import (
	"github.com/expandorg/requester-service/pkg/app"
	m "github.com/expandorg/requester-service/pkg/model"
	"github.com/globalsign/mgo/bson"
)

// Service for templates
type Service interface {
	List(ctx *app.ReqCtx) ([]*m.TaskTemplate, error)
	FindByID(ctx *app.ReqCtx, templateID bson.ObjectId) (*m.TaskTemplate, error)
}

type templateService struct {
	templates m.TemplateRepository
}

// New returns a new instance of a temapltes service.
func NewService(t m.TemplateRepository) Service {
	return &templateService{
		templates: t,
	}
}

func (r *templateService) List(ctx *app.ReqCtx) ([]*m.TaskTemplate, error) {
	return r.templates.List(ctx)
}

func (r *templateService) FindByID(ctx *app.ReqCtx, templateID bson.ObjectId) (*m.TaskTemplate, error) {
	return r.templates.FindByID(ctx, templateID)
}
