package draftservice

import (
	"github.com/expandorg/requester-service/pkg/app"
	m "github.com/expandorg/requester-service/pkg/model"
	"github.com/expandorg/requester-service/pkg/svc-kit/mongo"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type draftRepository struct {
	mongo.Repository
}

// NewRepository returns a new instance of a MongoDB draft repository.
func NewRepository() m.DraftRepository {
	return &draftRepository{
		Repository: mongo.Repository{Name: "drafts"},
	}
}

func (r *draftRepository) List(ctx *app.ReqCtx) ([]*m.Draft, error) {
	drafts := []*m.Draft{}
	err := r.Col(ctx).Find(nil).All(&drafts)
	return drafts, err
}

func (r *draftRepository) ListByUserID(ctx *app.ReqCtx, userID uint64) ([]*m.Draft, error) {
	drafts := []*m.Draft{}
	err := r.findCriteria(ctx, bson.M{"requesterId": userID}).All(&drafts)
	return drafts, err
}

func (r *draftRepository) ListByStatus(ctx *app.ReqCtx, status m.StatusType) ([]*m.Draft, error) {
	drafts := []*m.Draft{}
	err := r.findCriteria(ctx, bson.M{"status": status}).All(&drafts)
	return drafts, err
}

func (r *draftRepository) ListByUserIDAndStatus(ctx *app.ReqCtx, userID uint64, status m.StatusType) ([]*m.Draft, error) {
	drafts := []*m.Draft{}
	err := r.findCriteria(ctx, bson.M{"requesterId": userID, "status": status}).All(&drafts)
	return drafts, err
}

func (r *draftRepository) findCriteria(ctx *app.ReqCtx, query interface{}) *mgo.Query {
	return r.Col(ctx).Find(query)
}

func (r *draftRepository) FindByID(ctx *app.ReqCtx, draftID bson.ObjectId) (*m.Draft, error) {
	draft := new(m.Draft)
	err := r.Col(ctx).FindId(draftID).One(draft)
	return draft, err
}

func (r *draftRepository) Insert(ctx *app.ReqCtx, draft *m.Draft) error {
	return r.Col(ctx).Insert(draft)
}

func (r *draftRepository) Update(ctx *app.ReqCtx, draft *m.Draft) error {
	return r.Col(ctx).UpdateId(draft.ID, draft)
}

func (r *draftRepository) Delete(ctx *app.ReqCtx, draftID bson.ObjectId) error {
	return r.Col(ctx).RemoveId(draftID)
}
