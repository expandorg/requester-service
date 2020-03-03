package endpoints

import (
	"github.com/expandorg/requester-service/pkg/draftservice"
	"github.com/expandorg/requester-service/pkg/drafttemplates"
	"github.com/expandorg/requester-service/pkg/http/app/ctx"
	"github.com/expandorg/requester-service/pkg/http/app/permissions"
	"github.com/expandorg/requester-service/pkg/publisher"
	"github.com/expandorg/requester-service/pkg/svc-kit/http/api"
	"github.com/expandorg/requester-service/pkg/taskdata"
	"github.com/gin-gonic/gin"
)

// DraftRoutes setup templates routes
func DraftRoutes(r *gin.RouterGroup, h *api.Handler, drafts draftservice.Service, templates drafttemplates.Service, data taskdata.Service, publisher publisher.Service) {
	c := draftsController{
		drafts:    drafts,
		templates: templates,
		data:      data,
		publisher: publisher,
	}

	r.POST("/drafts", h.Handle(c.Create))

	perm := permissions.SetOwnDraft(h, drafts)

	r.GET("/drafts/:draftID", perm, h.Handle(c.Get))
	r.POST("/drafts/:draftID/variables", perm, h.Handle(c.UpdateVariables))
	r.POST("/drafts/:draftID/settings", perm, h.Handle(c.UpdateSettings))
	r.POST("/drafts/:draftID/verification", perm, h.Handle(c.UpdateVerification))

	r.POST("/drafts/:draftID/task/form", perm, h.Handle(c.Update))
	r.POST("/drafts/:draftID/verification/form", perm, h.Handle(c.Update))
	r.POST("/drafts/:draftID/onboarding", perm, h.Handle(c.Update))
	r.POST("/drafts/:draftID/funding", perm, h.Handle(c.Update))

	// r.POST("/drafts/:draftID/whitelist", perm, h.Handle(c.UpdateWhitelist))

	r.POST("/drafts/:draftID/copy", perm, h.Handle(c.Copy))
	r.DELETE("/drafts/:draftID", perm, h.Handle(c.Delete))
	r.POST("/drafts/:draftID/prepublish", perm, h.Handle(c.Prepublish))
}

type draftsController struct {
	drafts    draftservice.Service
	templates drafttemplates.Service
	data      taskdata.Service
	publisher publisher.Service
}

func (d *draftsController) Get(c *gin.Context) (interface{}, error) {
	return gin.H{"draft": ctx.Draft(c)}, nil
}

func (d *draftsController) Create(c *gin.Context) (interface{}, error) {
	body, err := ctx.BindCreateDraftRequest(c)
	if err != nil {
		return nil, api.BadRequest(err)
	}

	draft, err := d.drafts.Create(ctx.ReqCtx(c), ctx.UserID(c), body.TemplateID)
	if err != nil {
		return nil, api.ErrWithBody(err, body)
	}
	return gin.H{"draft": draft}, nil
}

func (d *draftsController) Delete(c *gin.Context) (interface{}, error) {
	draft := ctx.Draft(c)
	err := d.drafts.Delete(ctx.ReqCtx(c), draft)
	if err != nil {
		return nil, err
	}

	return gin.H{"draftId": draft.ID}, nil
}

func (d *draftsController) UpdateSettings(c *gin.Context) (interface{}, error) {
	body, err := ctx.BindUpdateDraftSettingsRequest(c)
	if err != nil {
		return nil, api.BadRequest(err)
	}

	draft, err := d.drafts.UpdateSettings(ctx.ReqCtx(c), ctx.Draft(c), body)
	if err != nil {
		return nil, api.ErrWithBody(err, body)
	}
	return gin.H{"draft": draft}, nil
}

func (d *draftsController) UpdateVariables(c *gin.Context) (interface{}, error) {
	body, err := ctx.BindUpdateDraftVariablesRequest(c)
	if err != nil {
		return nil, api.BadRequest(err)
	}

	draft, err := d.drafts.UpdateVariables(ctx.ReqCtx(c), ctx.Draft(c), body.Variables)
	if err != nil {
		return nil, api.ErrWithBody(err, body)
	}
	return gin.H{"draft": draft}, nil
}

func (d *draftsController) UpdateVerification(c *gin.Context) (interface{}, error) {
	body, err := ctx.BindUpdateDraftVerificationRequest(c)
	if err != nil {
		return nil, api.BadRequest(err)
	}

	draft, err := d.drafts.UpdateVerification(ctx.ReqCtx(c), ctx.Draft(c), body)
	if err != nil {
		return nil, api.ErrWithBody(err, body)
	}
	return gin.H{"draft": draft}, nil
}

func (d *draftsController) Update(c *gin.Context) (interface{}, error) {
	body, err := ctx.BindUpdateDraftRequest(c)
	if err != nil {
		return nil, api.BadRequest(err)
	}

	draft, err := d.drafts.Update(ctx.ReqCtx(c), ctx.Draft(c), body)
	if err != nil {
		return nil, api.ErrWithBody(err, body)
	}
	return gin.H{"draft": draft}, nil
}

func (d *draftsController) Copy(c *gin.Context) (interface{}, error) {
	draft, err := d.drafts.Copy(ctx.ReqCtx(c), ctx.Draft(c))
	if err != nil {
		return nil, err
	}
	return gin.H{"draft": draft}, nil
}

func (d *draftsController) Prepublish(c *gin.Context) (interface{}, error) {
	draft, err := d.publisher.Prepublish(ctx.ReqCtx(c), ctx.Draft(c), ctx.JWT(c))
	if err != nil {
		return nil, err
	}
	return gin.H{"draft": draft}, nil
}
