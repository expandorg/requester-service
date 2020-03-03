package drafttemplates

import (
	"github.com/globalsign/mgo/bson"

	"github.com/expandorg/requester-service/pkg/app"
	m "github.com/expandorg/requester-service/pkg/model"
	"github.com/expandorg/requester-service/pkg/svc-kit/mongo"
)

type templateRepository struct {
	mongo.Repository
}

// NewRepository returns a new instance of a templates repository
func NewRepository() m.TemplateRepository {
	return &templateRepository{
		Repository: mongo.Repository{Name: "taskTemplates"},
	}
}

func (r *templateRepository) List(ctx *app.ReqCtx) ([]*m.TaskTemplate, error) {
	templates := []*m.TaskTemplate{}
	err := r.Col(ctx).Find(nil).Sort("order").All(&templates)
	return templates, err
}

func (r *templateRepository) FindByID(ctx *app.ReqCtx, templateID bson.ObjectId) (*m.TaskTemplate, error) {
	template := new(m.TaskTemplate)
	err := r.Col(ctx).FindId(templateID).One(template)
	return template, err
}
