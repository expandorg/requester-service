package mongo

import (
	"github.com/expandorg/requester-service/pkg/app"
	"github.com/globalsign/mgo"
)

// Mongo  Repository base struct.
type Repository struct {
	Name string
}

// Col - returns mongo collection
func (r *Repository) Col(ctx *app.ReqCtx) *mgo.Collection {
	return ctx.Db().C(r.Name)
}
