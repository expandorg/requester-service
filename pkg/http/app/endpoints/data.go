package endpoints

import (
	"github.com/expandorg/requester-service/pkg/dataupload"
	"github.com/expandorg/requester-service/pkg/draftservice"
	"github.com/expandorg/requester-service/pkg/http/app/ctx"
	"github.com/expandorg/requester-service/pkg/http/app/permissions"
	m "github.com/expandorg/requester-service/pkg/model"
	"github.com/expandorg/requester-service/pkg/taskdata"
	"github.com/gemsorg/svc-kit/http/api"
	"github.com/gin-gonic/gin"
)

// DataRoutes setup templates routes
func DataRoutes(r *gin.RouterGroup, h *api.Handler, drafts draftservice.Service, data taskdata.Service, uploader dataupload.Service) {
	c := dataController{
		drafts:   drafts,
		data:     data,
		uploader: uploader,
	}

	ownDraft := permissions.SetOwnDraft(h, drafts)
	ownData := permissions.SetOwnData(h, data)

	r.GET("/drafts/:draftID/data/:dataID", ownDraft, ownData, h.Handle(c.Get))
	r.POST("/drafts/:draftID/data", ownDraft, h.Handle(c.Create))
	r.POST("/drafts/:draftID/data/:dataID/columns", ownDraft, ownData, h.Handle(c.UpdateColumns))
	r.DELETE("/drafts/:draftID/data", ownDraft, h.Handle(c.Delete))
}

type dataController struct {
	drafts   draftservice.Service
	data     taskdata.Service
	uploader dataupload.Service
}

func (dc *dataController) Get(c *gin.Context) (interface{}, error) {
	page := ctx.Page(c)
	data, pages, err := dc.data.FindPaged(ctx.ReqCtx(c), ctx.DraftID(c), ctx.DataID(c), page)
	if err != nil {
		return nil, err
	}
	return gin.H{"data": data, "pagination": m.NewPagination(page, pages)}, nil
}

func (dc *dataController) Create(c *gin.Context) (interface{}, error) {
	file, err := ctx.FileHeader(c, "data")
	if err != nil {
		return nil, api.BadRequest(err)
	}
	draft, err := dc.uploader.Upload(ctx.ReqCtx(c), ctx.UserID(c), ctx.Draft(c), file)
	if err != nil {
		return nil, err
	}
	return gin.H{"draft": draft}, nil
}

func (dc *dataController) UpdateColumns(c *gin.Context) (interface{}, error) {
	body, err := ctx.BindDataColumnsReq(c)
	if err != nil {
		return nil, api.BadRequest(err)
	}

	data, err := dc.data.UpdateColumns(ctx.ReqCtx(c), ctx.TaskData(c), body)
	if err != nil {
		return nil, api.ErrWithBody(err, body)
	}
	return gin.H{"data": data}, nil
}

func (dc *dataController) Delete(c *gin.Context) (interface{}, error) {
	draft, err := dc.drafts.DeleteTaskData(ctx.ReqCtx(c), ctx.Draft(c))
	if err != nil {
		return nil, err
	}
	return gin.H{"draft": draft}, nil
}
