package taskdata

import (
	"github.com/expandorg/requester-service/pkg/app"
	m "github.com/expandorg/requester-service/pkg/model"
	"github.com/globalsign/mgo/bson"
)

type Service interface {
	Find(ctx *app.ReqCtx, draftID bson.ObjectId, dataID bson.ObjectId) (*m.TaskData, error)
	FindPaged(ctx *app.ReqCtx, draftID bson.ObjectId, dataID bson.ObjectId, page uint64) (*m.TaskData, uint64, error)

	UpdateColumns(ctx *app.ReqCtx, taskData *m.TaskData, columnsReq *ColumnsRequest) (*m.TaskData, error)
	Delete(ctx *app.ReqCtx, ID bson.ObjectId) error
}

type dataService struct {
	data m.DataRepository
}

// New returns a new instance of a data service
func NewService(data m.DataRepository) Service {
	return &dataService{
		data: data,
	}
}

func (r *dataService) Find(ctx *app.ReqCtx, draftID bson.ObjectId, dataID bson.ObjectId) (*m.TaskData, error) {
	return r.data.Find(ctx, draftID, dataID)
}

func (r *dataService) FindPaged(ctx *app.ReqCtx, draftID bson.ObjectId, dataID bson.ObjectId, page uint64) (*m.TaskData, uint64, error) {
	return r.data.FindPaged(ctx, draftID, dataID, page)
}

func (r *dataService) UpdateColumns(ctx *app.ReqCtx, taskData *m.TaskData, c *ColumnsRequest) (*m.TaskData, error) {
	err := r.data.UpdateColumns(ctx, taskData, c.Columns)
	if err != nil {
		return nil, err
	}
	return r.data.Find(ctx, taskData.DraftID, taskData.ID)
}

func (r *dataService) Delete(ctx *app.ReqCtx, ID bson.ObjectId) error {
	return r.data.Delete(ctx, ID)
}

type ColumnsRequest struct {
	Columns []*m.TaskDataColumn `json:"columns"`
}
