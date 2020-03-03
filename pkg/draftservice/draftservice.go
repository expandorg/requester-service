package draftservice

import (
	"github.com/expandorg/requester-service/pkg/app"
	m "github.com/expandorg/requester-service/pkg/model"
	"github.com/expandorg/requester-service/pkg/svc-kit/svc"
	"github.com/globalsign/mgo/bson"
)

type draftService struct {
	drafts    m.DraftRepository
	data      m.DataRepository
	templates m.TemplateRepository
}

// NewService returns a new instance of a drafts service
func NewService(drafts m.DraftRepository, data m.DataRepository, templates m.TemplateRepository) Service {
	return &draftService{
		drafts:    drafts,
		data:      data,
		templates: templates,
	}
}

func (s *draftService) GetDraft(ctx *app.ReqCtx, draftID bson.ObjectId) (*m.Draft, error) {
	return s.drafts.FindByID(ctx, draftID)
}

func (s *draftService) Create(ctx *app.ReqCtx, userID uint64, templateID bson.ObjectId) (*m.Draft, error) {
	template, err := s.templates.FindByID(ctx, templateID)

	if err != nil {
		return nil, svc.NotFound(err)
	}

	draft := &m.Draft{
		ID:               bson.NewObjectId(),
		RequesterID:      userID,
		Status:           m.DraftStatus,
		TemplateID:       template.ID,
		Name:             template.Name,
		Onboarding:       template.Onboarding,
		TaskForm:         template.TaskForm,
		VerificationForm: template.VerificationForm,
		Eligibility:      template.Eligibility,
		Assignment:       template.Assignment,
		Verification:     template.Verification,
		Funding:          template.Funding,
		Variables:        template.Variables,
	}

	draft.Funding.Requirement = 0

	err = s.drafts.Insert(ctx, draft)
	if err != nil {
		return nil, err
	}
	return s.drafts.FindByID(ctx, draft.ID)
}

func (s *draftService) Copy(ctx *app.ReqCtx, draft *m.Draft) (*m.Draft, error) {
	draftID := bson.NewObjectId()

	if draft.DataID != "" {
		data, dataErr := s.data.Find(ctx, draft.ID, draft.DataID)
		if dataErr == nil {
			dataCopy, err := s.data.Copy(ctx, data, draftID)
			if err != nil {
				return nil, err
			}
			draft.DataID = dataCopy.ID
		}
	}
	draft.ID = draftID
	draft.Name = draft.Name + " copy"

	err := s.drafts.Insert(ctx, draft)
	if err != nil {
		return nil, err
	}
	return s.drafts.FindByID(ctx, draftID)
}

func (s *draftService) UpdateSettings(ctx *app.ReqCtx, draft *m.Draft, req *UpdateSettingsRequest) (*m.Draft, error) {
	draft.Name = req.Name
	draft.Description = req.Description
	draft.Logo = req.LogoURL

	draft.CallbackURL = req.CallbackURL

	if req.Staking {
		draft.Funding.Requirement = req.Stake
	} else {
		draft.Funding.Requirement = 0
	}

	err := s.drafts.Update(ctx, draft)
	if err != nil {
		return nil, err
	}
	return s.drafts.FindByID(ctx, draft.ID)
}

func (s *draftService) UpdateVariables(ctx *app.ReqCtx, draft *m.Draft, variables []string) (*m.Draft, error) {
	draft.Variables = variables
	err := s.drafts.Update(ctx, draft)
	if err != nil {
		return nil, err
	}
	return s.drafts.FindByID(ctx, draft.ID)
}

func (s *draftService) UpdateVerification(ctx *app.ReqCtx, draft *m.Draft, req *UpdateVerificationRequest) (*m.Draft, error) {

	if req.VerificationModule.IsValid() && draft.Verification.Module != req.VerificationModule {
		draft.Verification.Module = req.VerificationModule
	}
	if draft.Verification.Module == m.VerificationModuleConsensus {
		draft.Verification.AgreementCount = req.AgreementCount
	}

	draft.Verification.MinimumExecutionTime = req.MinimumExecutionTime

	err := s.drafts.Update(ctx, draft)
	if err != nil {
		return nil, err
	}
	return s.drafts.FindByID(ctx, draft.ID)
}

func (s *draftService) Update(ctx *app.ReqCtx, draft *m.Draft, req *UpdateRequest) (*m.Draft, error) {
	if req.Whitelist != nil {
		draft.Whitelist = req.Whitelist
	}

	if req.TaskForm != nil {
		draft.TaskForm = req.TaskForm
	}
	if req.VerificationForm != nil {
		draft.VerificationForm = req.VerificationForm
	}
	if req.Onboarding != nil {
		draft.Onboarding = req.Onboarding
	}

	if req.Eligibility != nil {
		draft.Eligibility = req.Eligibility
	}
	if req.Assignment != nil {
		draft.Assignment = req.Assignment
	}
	if req.Verification != nil {
		draft.Verification = req.Verification
	}
	if req.Funding != nil {
		draft.Funding = req.Funding
	}

	err := s.drafts.Update(ctx, draft)
	if err != nil {
		return nil, err
	}
	return s.drafts.FindByID(ctx, draft.ID)
}

func (s *draftService) Delete(ctx *app.ReqCtx, draft *m.Draft) error {
	err := s.deleteDataIfExist(ctx, draft.ID, draft.DataID)
	if err != nil {
		return err
	}
	return s.drafts.Delete(ctx, draft.ID)
}

func (s *draftService) deleteDataIfExist(ctx *app.ReqCtx, draftID bson.ObjectId, DataID bson.ObjectId) error {
	// To validate the DataID we will retrieve it and see if the data is blank
	// To Do: Check the length of DataID
	_, err := s.data.Find(ctx, draftID, DataID)
	if err != nil {
		return nil
	}
	return s.data.Delete(ctx, DataID)
}

func (s *draftService) DeleteTaskData(ctx *app.ReqCtx, draft *m.Draft) (*m.Draft, error) {
	err := s.data.Delete(ctx, draft.DataID)
	if err != nil {
		return nil, err
	}
	draft.DataID = ""
	err = s.drafts.Update(ctx, draft)
	if err != nil {
		return nil, err
	}
	return s.drafts.FindByID(ctx, draft.ID)
}
