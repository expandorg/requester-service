package model

import (
	"github.com/expandorg/requester-service/pkg/app"
	"github.com/globalsign/mgo/bson"
)

type TaskTemplate struct {
	ID               bson.ObjectId      `bson:"_id" json:"id"`
	Order            int32              `bson:"order" json:"order"`
	Name             string             `bson:"name" json:"name"`
	Onboarding       *DraftOnboarding   `bson:"onboarding" json:"onboarding"`
	Variables        []string           `bson:"variables,omitempty" json:"variables,omitempty"`
	DataSample       map[string]string  `bson:"dataSample,omitempty" json:"dataSample,omitempty"`
	TaskForm         *Form              `bson:"taskForm" json:"taskForm"`
	VerificationForm *Form              `bson:"verificationForm" json:"verificationForm"`
	Eligibility      *DraftEligibility  `bson:"eligibility,omitempty" json:"eligibility,omitempty"`
	Assignment       *DraftAssignment   `bson:"assignment,omitempty" json:"assignment,omitempty"`
	Verification     *DraftVerification `bson:"verification,omitempty" json:"verification,omitempty"`
	Funding          *DraftFunding      `bson:"funding,omitempty" json:"funding,omitempty"`
}

// TemplateRepository provide access to stored template entities.
type TemplateRepository interface {
	List(ctx *app.ReqCtx) ([]*TaskTemplate, error)
	FindByID(ctx *app.ReqCtx, templateID bson.ObjectId) (*TaskTemplate, error)
}
