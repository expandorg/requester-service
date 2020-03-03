package onboardingtemplate

import (
	"github.com/globalsign/mgo/bson"

	"github.com/expandorg/requester-service/pkg/app"
	m "github.com/expandorg/requester-service/pkg/model"
	"github.com/expandorg/requester-service/pkg/svc-kit/mongo"
)

type repository struct {
	mongo.Repository
}

// NewRepository returns a new instance of a onboarding templates repository
func NewRepository() m.OnboardingTemplateRepository {
	return &repository{
		Repository: mongo.Repository{Name: "onboardingTemplates"},
	}
}

func (r *repository) List(ctx *app.ReqCtx) ([]*m.OnboardingTemplate, error) {
	templates := []*m.OnboardingTemplate{}
	err := r.Col(ctx).Find(nil).All(&templates)
	return templates, err
}

func (r *repository) FindByID(ctx *app.ReqCtx, templateID bson.ObjectId) (*m.OnboardingTemplate, error) {
	template := new(m.OnboardingTemplate)
	err := r.Col(ctx).FindId(templateID).One(template)
	return template, err
}
