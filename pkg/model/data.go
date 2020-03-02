package model

import (
	"github.com/expandorg/requester-service/pkg/app"
	"github.com/expandorg/requester-service/pkg/upload"
	"github.com/globalsign/mgo/bson"
)

// DataRepository provide access to stored template entities.
type DataRepository interface {
	Find(ctx *app.ReqCtx, draftID bson.ObjectId, dataID bson.ObjectId) (*TaskData, error)
	FindPaged(ctx *app.ReqCtx, draftID bson.ObjectId, dataID bson.ObjectId, page uint64) (*TaskData, uint64, error)

	Create(ctx *app.ReqCtx, draftID bson.ObjectId, d *upload.Data) (bson.ObjectId, error)
	Copy(ctx *app.ReqCtx, data *TaskData, draftID bson.ObjectId) (*TaskData, error)
	UpdateColumns(ctx *app.ReqCtx, taskData *TaskData, columns []*TaskDataColumn) error
	Delete(ctx *app.ReqCtx, ID bson.ObjectId) error
}

type TaskData struct {
	ID      bson.ObjectId     `bson:"_id" json:"id"`
	DraftID bson.ObjectId     `bson:"draftId" json:"draftId"`
	Columns []*TaskDataColumn `bson:"columns" json:"columns"`
	Values  [][]string        `bson:"values" json:"values"`
}

type TaskDataColumn struct {
	Name     string `bson:"name" json:"name"`
	Variable string `bson:"variable" json:"variable"`
	Skipped  bool   `bson:"skipped" json:"skipped"`
}
