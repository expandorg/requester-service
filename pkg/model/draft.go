package model

import (
	"time"

	"github.com/expandorg/requester-service/pkg/app"
	"github.com/globalsign/mgo/bson"
)

// DraftRepository provide access to stored repository entities.
type DraftRepository interface {
	List(ctx *app.ReqCtx) ([]*Draft, error)
	ListByUserID(ctx *app.ReqCtx, userID uint64) ([]*Draft, error)
	ListByStatus(ctx *app.ReqCtx, status StatusType) ([]*Draft, error)
	ListByUserIDAndStatus(ctx *app.ReqCtx, userID uint64, status StatusType) ([]*Draft, error)

	FindByID(ctx *app.ReqCtx, draftID bson.ObjectId) (*Draft, error)

	Insert(ctx *app.ReqCtx, draft *Draft) error
	Update(ctx *app.ReqCtx, draft *Draft) error
	Delete(ctx *app.ReqCtx, draftID bson.ObjectId) error
}

type StatusType string

const (
	DraftStatus      StatusType = "draft"
	InProgressStatus StatusType = "in-progress"
	CompletedStatus  StatusType = "completed"
	PendingStatus    StatusType = "pending"
)

type EligibilityModuleID string

const (
	EligibilityModuleNoop  EligibilityModuleID = "noop"
	EligibilityModuleAll   EligibilityModuleID = "all"
	EligibilityModuleEmail EligibilityModuleID = "email"
)

type AssignmentModuleID string

const (
	AssignmentModuleNoop     AssignmentModuleID = "noop"
	AssignmentModuleOne      AssignmentModuleID = "one"
	AssignmentModuleAll      AssignmentModuleID = "all"
	AssignmentModuleExternal AssignmentModuleID = "external"
)

type FundingModuleID string

const (
	FundingModuleNoop        FundingModuleID = "noop"
	FundingModuleRequirement FundingModuleID = "requirement"
	FundingModuleStake       FundingModuleID = "stake"
)

type VerificationModuleID string

const (
	VerificationModuleNoop           VerificationModuleID = "noop"
	VerificationModuleRequester      VerificationModuleID = "requester"
	VerificationModuleConsensus      VerificationModuleID = "consensus"
	VerificationModuleAudit          VerificationModuleID = "audit"
	VerificationModuleAuditWhitelist VerificationModuleID = "audit-whitelist"
	VerificationModuleBulk           VerificationModuleID = "bulk"
)

func (id VerificationModuleID) IsValid() bool {
	return id == VerificationModuleNoop || id == VerificationModuleRequester || id == VerificationModuleConsensus || id == VerificationModuleAudit || id == VerificationModuleAuditWhitelist
}

func (status StatusType) IsValid() bool {
	return status == DraftStatus || status == InProgressStatus || status == CompletedStatus || status == PendingStatus
}

// Draft represent job draft entity
type Draft struct {
	ID               bson.ObjectId      `bson:"_id" json:"id"`
	Name             string             `bson:"name,omitempty" json:"name,omitempty"`
	Description      string             `bson:"description,omitempty" json:"description,omitempty"`
	Variables        []string           `bson:"variables,omitempty" json:"variables,omitempty"`
	Logo             string             `bson:"logo,omitempty" json:"logo,omitempty"`
	StartDate        *time.Time         `bson:"startDate,omitempty" json:"startDate,omitempty"`
	EndDate          *time.Time         `bson:"endDate,omitempty" json:"endDate,omitempty"`
	RequesterID      uint64             `bson:"requesterId" json:"requesterId"`
	JobID            uint64             `bson:"jobId,omitempty" json:"jobId,omitempty"`
	Status           StatusType         `bson:"status" json:"status"`
	DataID           bson.ObjectId      `bson:"dataId,omitempty" json:"dataId,omitempty"`
	TemplateID       bson.ObjectId      `bson:"templateId,omitempty" json:"templateId,omitempty"`
	TaskForm         *Form              `bson:"taskForm,omitempty" json:"taskForm,omitempty"`
	VerificationForm *Form              `bson:"verificationForm,omitempty" json:"verificationForm,omitempty"`
	Onboarding       *DraftOnboarding   `bson:"onboarding,omitempty" json:"onboarding,omitempty"`
	Eligibility      *DraftEligibility  `bson:"eligibility,omitempty" json:"eligibility,omitempty"`
	Assignment       *DraftAssignment   `bson:"assignment,omitempty" json:"assignment,omitempty"`
	Verification     *DraftVerification `bson:"verification,omitempty" json:"verification,omitempty"`
	Funding          *DraftFunding      `bson:"funding,omitempty" json:"funding,omitempty"`
	Whitelist        []*DraftWhitelist  `bson:"whitelist,omitempty" json:"whitelist,omitempty"`
	CallbackURL      string             `bson:"callbackUrl,omitempty" json:"callbackUrl,omitempty"`
}

// type DraftLogic struct {
// 	Eligibility  *DraftEligibility  `bson:"eligibility,omitempty" json:"eligibility,omitempty"`
// 	Assignment   *DraftAssignment   `bson:"assignment,omitempty" json:"assignment,omitempty"`
// 	Verification *DraftVerification `bson:"verification,omitempty" json:"verification,omitempty"`
// 	Funding      *DraftFunding      `bson:"funding,omitempty" json:"funding,omitempty"`
// }

type DraftOnboarding struct {
	Enabled bool                  `bson:"enabled" json:"enabled"`
	Steps   []DraftOnboardingStep `bson:"steps,omitempty" json:"steps,omitempty"`
}

type DraftOnboardingStep struct {
	ID             string               `bson:"_id" json:"id"`
	Name           string               `bson:"name" json:"name"`
	Form           *Form                `bson:"form" json:"form"`
	IsGroup        bool                 `bson:"isGroup,omitempty" json:"isGroup,omitempty"`
	ScoreThreshold int64                `bson:"scoreThreshold,omitempty" json:"scoreThreshold,omitempty"`
	Retries        int64                `bson:"retries,omitempty" json:"retries,omitempty"`
	FailureMessage string               `bson:"failureMessage,omitempty" json:"failureMessage,omitempty"`
	Data           *OnboardingGroupData `bson:"data,omitempty" json:"data,omitempty"`
}

type DraftEligibility struct {
	Module EligibilityModuleID `bson:"module,omitempty" json:"module,omitempty"`
}

type DraftAssignment struct {
	Module     AssignmentModuleID `bson:"module,omitempty" json:"module,omitempty"`
	Limit      int64              `bson:"limit,omitempty" json:"limit,omitempty"`
	Repeat     bool               `bson:"repeat,omitempty" json:"repeat,omitempty"`
	Expiration int64              `bson:"expiration,omitempty" json:"expiration,omitempty"`
}

type DraftVerification struct {
	Module               VerificationModuleID `bson:"module,omitempty" json:"module,omitempty"`
	AgreementCount       int64                `bson:"agreementCount,omitempty" json:"agreementCount,omitempty"`
	ScoreThreshold       float64              `bson:"scoreThreshold,omitempty" json:"scoreThreshold,omitempty"`
	MinimumExecutionTime uint64               `bson:"minimumExecutionTime,omitempty" json:"minimumExecutionTime,omitempty"`
}

type DraftFunding struct {
	Module             FundingModuleID `bson:"module,omitempty" json:"module,omitempty"`
	Requirement        float64         `bson:"requirement,omitempty" json:"requirement,omitempty"`
	Balance            float64         `bson:"balance,omitempty" json:"balance,omitempty"`
	Reward             float64         `bson:"reward,omitempty" json:"reward,omitempty"`
	VerificationReward float64         `bson:"verificationReward,omitempty" json:"verificationReward,omitempty"`
}

type DraftWhitelist struct {
	ID    bson.ObjectId `bson:"_id" json:"id"`
	Param string        `bson:"param" json:"param"`
	Op    string        `bson:"op" json:"op"`
	Value string        `bson:"value" json:"value"`
}

type OnboardingGroupData struct {
	Answer  *OnboardingGroupDataAnswer     `bson:"answer" json:"answer"`
	Columns []*OnboardingGroupDataVariable `bson:"columns" json:"columns"`
	Steps   []*OnboardingGroupDataStep     `bson:"steps" json:"steps"`
}

type OnboardingGroupDataAnswer struct {
	Field string `bson:"field" json:"field"`
}

type OnboardingGroupDataVariable struct {
	Name string `bson:"name" json:"name"`
	Type string `bson:"type" json:"type"`
}

type OnboardingGroupDataStep struct {
	Answer string   `bson:"answer" json:"answer"`
	Values []string `bson:"values" json:"values"`
}
