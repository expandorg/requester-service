package endpoints

import (
	"github.com/expandorg/requester-service/pkg/dashboard"
	"github.com/expandorg/requester-service/pkg/draftservice"
	"github.com/expandorg/requester-service/pkg/http/app/ctx"
	"github.com/expandorg/requester-service/pkg/http/app/permissions"
	m "github.com/expandorg/requester-service/pkg/model"
	"github.com/expandorg/requester-service/pkg/publisher"
	"github.com/expandorg/requester-service/pkg/taskdata"
	"github.com/gemsorg/svc-kit/http/api"
	"github.com/gin-gonic/gin"
)

// AdminRoutes setup templates routes
func AdminRoutes(r *gin.RouterGroup, h *api.Handler, publisher publisher.Service, dashboard dashboard.Service, drafts draftservice.Service, data taskdata.Service, adminApiKey string) {
	c := adminController{
		publisher: publisher,
		dashboard: dashboard,
		drafts:    drafts,
		data:      data,
		apiKey:    adminApiKey,
	}
	isAdmin := permissions.IsAdmin(h)

	r.GET("/admin/review", isAdmin, h.Handle(c.GetAdminDrafts))
	r.GET("/admin/drafts/:draftID", isAdmin, h.Handle(c.GetDraft))
	r.GET("/admin/drafts/:draftID/data/:dataID", isAdmin, h.Handle(c.GetData))
	r.POST("/admin/drafts/:draftID/publish", isAdmin, h.Handle(c.PublishDraft))
	r.POST("/admin/drafts/:draftID/reject", isAdmin, h.Handle(c.RejectDraft))
}

type adminController struct {
	publisher publisher.Service
	dashboard dashboard.Service
	drafts    draftservice.Service
	data      taskdata.Service
	apiKey    string
}

func (ac *adminController) GetDraft(c *gin.Context) (interface{}, error) {
	draft, err := ac.drafts.GetDraft(ctx.ReqCtx(c), ctx.DraftID(c))
	if err != nil {
		return err, nil
	}
	return gin.H{"draft": draft}, nil
}

func (ac *adminController) PublishDraft(c *gin.Context) (interface{}, error) {
	draft, err := ac.publisher.AdminPublish(ctx.ReqCtx(c), ctx.DraftID(c), ac.apiKey)
	if err != nil {
		return nil, err
	}
	return gin.H{"draft": draft}, nil
}

func (ac *adminController) RejectDraft(c *gin.Context) (interface{}, error) {
	body, err := ctx.BindRejectDraftRequest(c)
	if err != nil {
		return nil, api.BadRequest(err)
	}

	err = ac.publisher.AdminReject(ctx.ReqCtx(c), ctx.DraftID(c), body.Message, ac.apiKey)
	if err != nil {
		return nil, err
	}
	return gin.H{"ok": true}, nil
}

func (ac *adminController) GetAdminDrafts(c *gin.Context) (interface{}, error) {
	drafts, err := ac.dashboard.GetAdminDrafts(ctx.ReqCtx(c))
	if err != nil {
		return nil, err
	}
	return gin.H{"drafts": drafts}, nil
}

func (ac *adminController) GetData(c *gin.Context) (interface{}, error) {
	page := ctx.Page(c)
	data, pages, err := ac.data.FindPaged(ctx.ReqCtx(c), ctx.DraftID(c), ctx.DataID(c), page)
	if err != nil {
		return nil, err
	}
	return gin.H{"data": data, "pagination": m.NewPagination(page, pages)}, nil
}
