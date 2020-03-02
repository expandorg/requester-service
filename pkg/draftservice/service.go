package draftservice

import (
	"github.com/expandorg/requester-service/pkg/app"

	m "github.com/expandorg/requester-service/pkg/model"
	"github.com/globalsign/mgo/bson"
)

// Service interface.
type Service interface {
	GetDraft(ctx *app.ReqCtx, draftID bson.ObjectId) (*m.Draft, error)

	Create(ctx *app.ReqCtx, userID uint64, templateID bson.ObjectId) (*m.Draft, error)

	UpdateSettings(ctx *app.ReqCtx, draft *m.Draft, req *UpdateSettingsRequest) (*m.Draft, error)
	UpdateVariables(ctx *app.ReqCtx, draft *m.Draft, variables []string) (*m.Draft, error)

	UpdateVerification(ctx *app.ReqCtx, draft *m.Draft, req *UpdateVerificationRequest) (*m.Draft, error)

	Update(ctx *app.ReqCtx, draft *m.Draft, req *UpdateRequest) (*m.Draft, error)

	Delete(ctx *app.ReqCtx, draft *m.Draft) error
	DeleteTaskData(ctx *app.ReqCtx, draft *m.Draft) (*m.Draft, error)

	Copy(ctx *app.ReqCtx, draft *m.Draft) (*m.Draft, error)
}

type UpdateSettingsRequest struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	LogoURL     string  `json:"logoUrl,omitempty"`
	Staking     bool    `json:"staking,omitempty"`
	Stake       float64 `json:"stake,omitempty"`
	Deduct      bool    `json:"deduct,omitempty"`
	CallbackURL string  `json:"callbackUrl,omitempty"`
}

type UpdateVariablesRequest struct {
	Variables []string `json:"variables"`
}

type UpdateVerificationRequest struct {
	VerificationModule   m.VerificationModuleID `json:"verificationModule"`
	AgreementCount       int64                  `json:"agreementCount,omitempty"`
	MinimumExecutionTime uint64                 `json:"minimumExecutionTime,omitempty"`
}

type UpdateRequest struct {
	Whitelist        []*m.DraftWhitelist  `json:"whitelist,omitempty"`
	TaskForm         *m.Form              `json:"taskForm,omitempty"`
	VerificationForm *m.Form              `json:"verificationForm,omitempty"`
	Onboarding       *m.DraftOnboarding   `json:"onboarding,omitempty"`
	Eligibility      *m.DraftEligibility  `json:"eligibility,omitempty"`
	Assignment       *m.DraftAssignment   `json:"assignment,omitempty"`
	Verification     *m.DraftVerification `json:"verification,omitempty"`
	Funding          *m.DraftFunding      `json:"funding,omitempty"`
}
type CreateRequest struct {
	TemplateID bson.ObjectId `json:"templateId"`
}
