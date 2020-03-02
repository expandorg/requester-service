package model

import (
	"github.com/expandorg/requester-service/pkg/app"
	"github.com/globalsign/mgo/bson"
)

// OnboardingTemplateRepository provide access to stored template entities.
type OnboardingTemplateRepository interface {
	List(ctx *app.ReqCtx) ([]*OnboardingTemplate, error)
	FindByID(ctx *app.ReqCtx, templateID bson.ObjectId) (*OnboardingTemplate, error)
}

// OnboardingTemplate for setting draft
type OnboardingTemplate struct {
	ID             bson.ObjectId        `bson:"_id" json:"id"`
	Name           string               `bson:"name" json:"name"`
	IsGroup        bool                 `bson:"isGroup,omitempty" json:"isGroup,omitempty"`
	ScoreThreshold int64                `bson:"scoreThreshold,omitempty" json:"scoreThreshold,omitempty"`
	Retries        int64                `bson:"retries,omitempty" json:"retries,omitempty"`
	FailureMessage string               `bson:"failureMessage,omitempty" json:"failureMessage,omitempty"`
	TaskForm       *Form                `bson:"taskForm" json:"taskForm"`
	Data           *OnboardingGroupData `bson:"data,omitempty" json:"data,omitempty"`
}
