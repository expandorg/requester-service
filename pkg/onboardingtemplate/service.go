package onboardingtemplate

import (
	"github.com/expandorg/requester-service/pkg/app"
	m "github.com/expandorg/requester-service/pkg/model"
	"github.com/globalsign/mgo/bson"
)

// Service provide access to stored template entities.
type Service interface {
	List(ctx *app.ReqCtx) ([]*m.OnboardingTemplate, error)
	FindByID(ctx *app.ReqCtx, templateID bson.ObjectId) (*m.OnboardingTemplate, error)
}

type service struct {
	templates m.OnboardingTemplateRepository
}

// NewService returns a new instance of an onboarding templates service
func NewService(o m.OnboardingTemplateRepository) Service {
	return &service{
		templates: o,
	}
}

func (r *service) List(ctx *app.ReqCtx) ([]*m.OnboardingTemplate, error) {
	return r.templates.List(ctx)
}

func (r *service) FindByID(ctx *app.ReqCtx, templateID bson.ObjectId) (*m.OnboardingTemplate, error) {
	return r.templates.FindByID(ctx, templateID)
}
